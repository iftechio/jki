package cp

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
	"github.com/iftechio/jki/pkg/factory"
	"github.com/iftechio/jki/pkg/image"
	"github.com/iftechio/jki/pkg/registry"
	"github.com/iftechio/jki/pkg/utils"
	"github.com/spf13/cobra"
)

type CopyOptions struct {
	resolver     *registry.Resolver
	dockerClient *client.Client
	dstRegistry  *registry.Registry
	saveImage    bool
}

func (o *CopyOptions) Complete(f factory.Factory, cmd *cobra.Command, args []string) error {
	var err error
	o.resolver, err = f.ToResolver()
	if err != nil {
		return err
	}

	o.dockerClient, err = f.DockerClient()
	if err != nil {
		return err
	}

	dstReg, registries, err := f.LoadRegistries()
	if err != nil {
		return err
	}
	if len(args) > 1 {
		dstReg = args[1]
	}
	if _, exist := registries[dstReg]; !exist {
		return fmt.Errorf("registry not found: %s", dstReg)
	}
	o.dstRegistry = registries[dstReg]
	return nil
}

func (o *CopyOptions) Validate(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("wrong number of arguments")
	}
	return nil
}

func (o *CopyOptions) Run(args []string) error {
	frImg := args[0]
	ctx := context.TODO()
	termFd, isTerm := term.GetFdInfo(os.Stdout)

	_, _, err := o.dockerClient.ImageInspectWithRaw(ctx, frImg)
	if err != nil {
		if client.IsErrNotFound(err) {
			// not exist locally
			// try to pull from registry
			reg, err := o.resolver.ResolveRegistryByImage(frImg)
			if err != nil {
				return err
			}
			frToken, err := reg.GetAuthToken()
			if err != nil {
				return err
			}
			frImg, err = o.completeImageStr(frImg, reg)
			if err != nil {
				return err
			}

			out, err := o.dockerClient.ImagePull(ctx, frImg, types.ImagePullOptions{RegistryAuth: frToken})
			if err != nil {
				return err
			}

			fmt.Printf("===== Pulling %s =====\n", frImg)
			defer out.Close()

			err = jsonmessage.DisplayJSONMessagesStream(out, os.Stdout, termFd, isTerm, nil)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	img := image.FromString(frImg)
	toReg := o.dstRegistry
	err = toReg.CreateRepoIfNotExists(img.Repo)
	if err != nil {
		return err
	}

	img.Domain = toReg.Prefix()
	toImg := img.String()
	_ = o.dockerClient.ImageTag(ctx, frImg, toImg)

	toToken, err := toReg.GetAuthToken()
	if err != nil {
		return err
	}
	pushOut, err := o.dockerClient.ImagePush(ctx, toImg, types.ImagePushOptions{RegistryAuth: toToken})
	if err != nil {
		return err
	}

	fmt.Printf("===== Pushing %s =====\n", toImg)
	defer pushOut.Close()

	err = jsonmessage.DisplayJSONMessagesStream(pushOut, os.Stdout, termFd, isTerm, nil)
	if err != nil {
		return err
	}

	if !o.saveImage {
		o.removeImages(ctx, frImg, toImg)
	}

	fmt.Println("Done!")
	return nil
}

func NewCopyOptions() *CopyOptions {
	return &CopyOptions{}
}

func NewCmdCp(f factory.Factory) *cobra.Command {
	o := NewCopyOptions()
	cmd := &cobra.Command{
		Use:   "cp <IMAGE> [REGISTRY NAME]",
		Short: "Copy images from one registry to another",
		Run: func(cmd *cobra.Command, args []string) {
			utils.CheckError(o.Complete(f, cmd, args))
			utils.CheckError(o.Validate(args))
			utils.CheckError(o.Run(args))
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&o.saveImage, "save-image", o.saveImage, "The local image will not be deleted after the copy is completed")
	return cmd
}

func (o *CopyOptions) removeImages(ctx context.Context, imageNames ...string) {
	for _, name := range imageNames {
		if _, err := o.dockerClient.ImageRemove(ctx, name, types.ImageRemoveOptions{}); err != nil {
			fmt.Printf("an error appears in removing image, err: %v\n", err)
		}
	}
	return
}

func (o *CopyOptions) completeImageStr(imgStr string, reg registry.RegistryInterface) (string, error) {
	if strings.Contains(imgStr, ":") {
		return imgStr, nil
	}
	splits := strings.Split(imgStr, "/")
	repo := splits[len(splits)-1]
	tag, err := reg.GetLatestTag(repo)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%s", imgStr, tag), nil
}
