package build

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
	cmdutils "github.com/bario/jki/pkg/cmd/utils"
	"github.com/bario/jki/pkg/git"
	"github.com/bario/jki/pkg/registry"
	"github.com/spf13/cobra"
)

func prompt(hint string) string {
	fmt.Print(hint)
	var input string
	fmt.Scanln(&input)
	return input
}

type BuildOptions struct {
	context        string
	dockerFileName string
	imageName      string
	registry       *registry.Registry
	dockerClient   *client.Client
}

func NewBuildOptions() *BuildOptions {
	return &BuildOptions{}
}

func (o *BuildOptions) Complete(f cmdutils.Factory, cmd *cobra.Command, args []string) error {
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
	o.registry = registries[defReg]
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
	currentHash, err := git.GetAbbrevCommitHash()
	if err != nil {
		return err
	}
	tag, err := git.GetTagOfCommit(currentHash)
	if err != nil {
		branch, err := git.GetCurrentBranch()
		if err != nil {
			return err
		}
		tag = fmt.Sprintf("%s-%s", branch, currentHash)
	}
	tag = strings.ReplaceAll(tag, "/", "-")
	tag = strings.ToLower(tag)

	domain := o.registry.Domain()
	ctx := context.TODO()
	tarStream, err := archive.TarWithOptions(o.context, &archive.TarOptions{})
	if err != nil {
		return fmt.Errorf("tar: %s", err)
	}
	defer tarStream.Close()
	image := fmt.Sprintf("%s/%s:%s", domain, o.imageName, tag)
	buildOpts := types.ImageBuildOptions{
		Tags:       []string{image},
		Remove:     true,
		Dockerfile: o.dockerFileName,
	}
	resp, err := o.dockerClient.ImageBuild(ctx, tarStream, buildOpts)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	termFd, isTerm := term.GetFdInfo(os.Stdout)
	err = jsonmessage.DisplayJSONMessagesStream(resp.Body, os.Stdout, termFd, isTerm, nil)
	if err != nil {
		return err
	}

	err = o.registry.CreateRepoIfNotExists(o.imageName)
	if err != nil {
		return err
	}

	authToken, err := o.registry.GetAuthToken()
	if err != nil {
		return err
	}
	pushResp, err := o.dockerClient.ImagePush(ctx, image, types.ImagePushOptions{RegistryAuth: authToken})
	if err != nil {
		return err
	}

	err = jsonmessage.DisplayJSONMessagesStream(pushResp, os.Stdout, termFd, isTerm, nil)
	if err != nil {
		return err
	}

	fmt.Println("Done!")
	return nil
}

func NewCmdBuild(f cmdutils.Factory) *cobra.Command {
	o := NewBuildOptions()
	cmd := &cobra.Command{
		Use:   "build [PATH]",
		Short: "Build docker image",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutils.CheckError(o.Complete(f, cmd, args))
			cmdutils.CheckError(o.Validate(args))
			cmdutils.CheckError(o.Run())
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
	return cmd
}
