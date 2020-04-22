package pull

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
	"github.com/spf13/cobra"

	"github.com/iftechio/jki/pkg/factory"
	"github.com/iftechio/jki/pkg/registry"
	"github.com/iftechio/jki/pkg/utils"
)

type Options struct {
	resolver     *registry.Resolver
	dockerClient *client.Client
	imageRef     string
}

func (o *Options) Complete(f factory.Factory, cmd *cobra.Command, args []string) error {
	var err error
	o.resolver, err = f.ToResolver()
	if err != nil {
		return err
	}

	o.dockerClient, err = f.DockerClient()
	if err != nil {
		return err
	}

	o.imageRef = args[0]
	return nil
}

func (o *Options) Validate(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("wrong number of arguments")
	}
	return nil
}

func (o *Options) Run() error {
	ctx := context.Background()
	termFd, isTerm := term.GetFdInfo(os.Stdout)

	_, _, err := o.dockerClient.ImageInspectWithRaw(ctx, o.imageRef)
	if err == nil {
		return nil
	} else if !client.IsErrNotFound(err) {
		return err
	}

	reg, err := o.resolver.ResolveRegistryByImage(o.imageRef)
	if err != nil {
		return err
	}
	token, err := reg.GetAuthToken()
	if err != nil {
		return err
	}
	out, err := o.dockerClient.ImagePull(ctx, o.imageRef, types.ImagePullOptions{RegistryAuth: token})
	if err != nil {
		return err
	}
	err = jsonmessage.DisplayJSONMessagesStream(out, os.Stdout, termFd, isTerm, nil)
	if err != nil {
		return err
	}
	return nil
}

func NewCmdPull(f factory.Factory) *cobra.Command {
	o := Options{}
	cmd := &cobra.Command{
		Use:   "pull <image>",
		Short: "Pull image from cloud registry",
		Run: func(cmd *cobra.Command, args []string) {
			utils.CheckError(o.Complete(f, cmd, args))
			utils.CheckError(o.Validate(args))
			utils.CheckError(o.Run())
		},
	}
	return cmd
}
