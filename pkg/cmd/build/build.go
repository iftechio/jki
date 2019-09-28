package build

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"unicode"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
	"github.com/bario/jki/pkg/factory"
	"github.com/bario/jki/pkg/git"
	"github.com/bario/jki/pkg/registry"
	"github.com/bario/jki/pkg/utils"
	"github.com/spf13/cobra"
)

func printInfo(msg string) {
	fmt.Printf(">>>>> %s\n", msg)
}

func prompt(hint string) string {
	fmt.Print(hint)
	var input string
	fmt.Scanln(&input)
	return input
}

func setClipboard(data string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xclip", "-selection", "c")
	case "darwin":
		cmd = exec.Command("pbcopy")
	default:
		return fmt.Errorf("%s not supported", runtime.GOOS)
	}
	cmd.Stdin = strings.NewReader(data)
	return cmd.Run()
}

func notifyUser(msg, title string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("notify-send", msg, title)
	case "darwin":
		cmd = exec.Command("osascript", "-e", fmt.Sprintf(`display notification "%s" with title "%s"`, msg, title))
	default:
		return fmt.Errorf("%s not supported", runtime.GOOS)
	}
	return cmd.Run()
}

type BuildOptions struct {
	context        string
	dockerFileName string
	imageName      string
	tagName        string
	dstRegistry    *registry.Registry
	allRegistries  map[string]*registry.Registry
	dockerClient   *client.Client
}

func NewBuildOptions() *BuildOptions {
	return &BuildOptions{}
}

func (o *BuildOptions) Complete(f factory.Factory, cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		o.context = args[0]
	}
	var err error
	o.dockerClient, err = f.DockerClient()
	if err != nil {
		return err
	}
	defReg, registries, err := f.LoadRegistries()
	if err != nil {
		return err
	}
	if strings.IndexFunc(o.imageName, unicode.IsUpper) != -1 {
		o.imageName = strings.ToLower(o.imageName)
		fmt.Printf("WARNING: uppercase char is not allowed in image name, changed to `%s`\n", o.imageName)
	}
	o.dstRegistry = registries[defReg]
	o.allRegistries = registries
	return nil
}

func (o *BuildOptions) Validate(args []string) error {
	return nil
}

func (o *BuildOptions) Run() error {
	if git.HasChanges() {
		input := strings.ToLower(prompt("当前有未提交的改动, 是否继续构建? (Y/n) "))
		if input == "n" {
			return nil
		}
	}

	var (
		tag string
		err error
	)
	if len(o.tagName) != 0 {
		tag = o.tagName
	} else {
		currentHash, err := git.GetAbbrevCommitHash()
		if err != nil {
			fmt.Println("WARNING: cannot get current commit, use `latest` as tag.")
			tag = "latest"
		} else {
			tag, err = git.GetTagOfCommit(currentHash)
			if err != nil {
				branch, err := git.GetCurrentBranch()
				if err != nil {
					return err
				}
				tag = fmt.Sprintf("%s-%s", branch, currentHash)
				tag = strings.ToLower(strings.ReplaceAll(tag, "/", "-"))
			}
		}
	}

	ctx := context.TODO()

	repoWithTag := fmt.Sprintf("%s:%s", o.imageName, tag)
	image := fmt.Sprintf("%s/%s", o.dstRegistry.Prefix(), repoWithTag)

	authConfigs := make(map[string]types.AuthConfig, len(o.allRegistries))
	dkfile, err := os.Open(o.dockerFileName)
	if err != nil {
		return err
	}
	defer dkfile.Close()
	baseImages, err := utils.ExtractBaseImages(dkfile)
	if err != nil {
		return err
	}
	mem := make(map[string]struct{}, len(o.allRegistries))
	for _, baseImage := range baseImages {
		for name, reg := range o.allRegistries {
			if _, ok := mem[name]; ok {
				continue
			}
			if strings.HasPrefix(baseImage, reg.Prefix()) {
				authCfg, err := reg.GetAuthConfig()
				if err != nil {
					return fmt.Errorf("get authconfig of %s: %s", name, err)
				}
				authConfigs[authCfg.ServerAddress] = authCfg
				mem[name] = struct{}{}
			}
		}
	}
	buildOpts := types.ImageBuildOptions{
		Tags:        []string{image},
		Remove:      true,
		Dockerfile:  o.dockerFileName,
		AuthConfigs: authConfigs,
	}

	ignores, err := utils.ReadDockerIgnore(o.context)
	if err != nil {
		return err
	}
	tarStream, err := archive.TarWithOptions(o.context, &archive.TarOptions{
		ExcludePatterns: ignores,
	})
	if err != nil {
		return fmt.Errorf("tar: %s", err)
	}
	defer tarStream.Close()

	resp, err := o.dockerClient.ImageBuild(ctx, tarStream, buildOpts)
	if err != nil {
		_ = notifyUser(" ", "镜像构建失败")
		return err
	}
	printInfo("开始构建镜像")
	defer resp.Body.Close()
	termFd, isTerm := term.GetFdInfo(os.Stdout)
	err = jsonmessage.DisplayJSONMessagesStream(resp.Body, os.Stdout, termFd, isTerm, nil)
	if err != nil {
		_ = notifyUser(" ", "镜像构建失败")
		return err
	}
	printInfo("镜像构建成功")

	err = o.dstRegistry.CreateRepoIfNotExists(o.imageName)
	if err != nil {
		return err
	}

	authToken, err := o.dstRegistry.GetAuthToken()
	if err != nil {
		return err
	}

	pushResp, err := o.dockerClient.ImagePush(ctx, image, types.ImagePushOptions{RegistryAuth: authToken})
	if err != nil {
		_ = notifyUser(" ", "镜像上传失败")
		return err
	}

	printInfo("开始上传镜像")
	defer pushResp.Close()
	err = jsonmessage.DisplayJSONMessagesStream(pushResp, os.Stdout, termFd, isTerm, nil)
	if err != nil {
		_ = notifyUser(" ", "镜像上传失败")
		return err
	}

	fmt.Println("镜像上传成功:")
	fmt.Println(image)
	_ = setClipboard(image)
	printInfo("镜像地址已复制到粘贴板")
	_ = notifyUser(repoWithTag, "镜像构建并上传成功")
	return nil
}

func NewCmdBuild(f factory.Factory) *cobra.Command {
	o := NewBuildOptions()
	cmd := &cobra.Command{
		Use:   "build [PATH]",
		Short: "Build docker image",
		Run: func(cmd *cobra.Command, args []string) {
			utils.CheckError(o.Complete(f, cmd, args))
			utils.CheckError(o.Validate(args))
			utils.CheckError(o.Run())
		},
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	o.context = wd

	flags := cmd.Flags()
	flags.StringVarP(&o.dockerFileName, "file", "f", "Dockerfile", "Name of the Dockerfile")
	flags.StringVar(&o.imageName, "image-name", path.Base(wd), "Custom image name")
	flags.StringVarP(&o.tagName, "tag-name", "t", "", "Custom tag name")
	return cmd
}
