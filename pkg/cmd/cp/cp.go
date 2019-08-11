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
	cmdutils "github.com/bario/jki/pkg/cmd/utils"
	"github.com/bario/jki/pkg/registry"
	"github.com/spf13/cobra"
)

type CopyOptions struct {
	resolver     *registry.Resolver
	dockerClient *client.Client
	dstRegistry  *registry.Registry
}

func (o *CopyOptions) Complete(f cmdutils.Factory, cmd *cobra.Command, args []string) error {
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

	frToken, err := o.resolver.ResolveRegistryAuth(frImg)
	if err != nil {
		return err
	}

	out, err := o.dockerClient.ImagePull(ctx, frImg, types.ImagePullOptions{RegistryAuth: frToken})
	if err != nil {
		return err
	}

	fmt.Printf("===== Pulling %s =====\n", frImg)
	defer out.Close()

	termFd, isTerm := term.GetFdInfo(os.Stdout)

	err = jsonmessage.DisplayJSONMessagesStream(out, os.Stdout, termFd, isTerm, nil)
	if err != nil {
		return err
	}

	var toRepo string
	parts := strings.Split(frImg, "/")
	repoWithTag := parts[len(parts)-1]
	parts = strings.Split(repoWithTag, ":")
	if len(parts) == 1 {
		// missing colon
		repoWithTag += ":latest"
		frImg += ":latest"
		toRepo = parts[0]
	} else {
		toRepo = parts[0]
	}

	toReg := o.dstRegistry

	err = toReg.CreateRepoIfNotExists(toRepo)
	if err != nil {
		return err
	}

	domain := toReg.Domain()
	toImg := domain + "/" + repoWithTag
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
	fmt.Println("Done!")
	return nil
}

func NewCopyOptions() *CopyOptions {
	return &CopyOptions{}
}

func NewCmdCp(f cmdutils.Factory) *cobra.Command {
	o := NewCopyOptions()
	cmd := &cobra.Command{
		Use:   "cp <IMAGE> [REGISTRY NAME]",
		Short: "Copy images from one registry to another",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutils.CheckError(o.Complete(f, cmd, args))
			cmdutils.CheckError(o.Validate(args))
			cmdutils.CheckError(o.Run(args))
		},
	}
	return cmd
}
