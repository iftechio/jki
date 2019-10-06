package main

import (
	"log"

	"github.com/iftechio/jki/pkg/cmd/build"
	"github.com/iftechio/jki/pkg/cmd/completion"
	"github.com/iftechio/jki/pkg/cmd/config"
	"github.com/iftechio/jki/pkg/cmd/cp"
	"github.com/iftechio/jki/pkg/cmd/transferimage"
	"github.com/iftechio/jki/pkg/cmd/version"
	"github.com/iftechio/jki/pkg/configflags"
	"github.com/iftechio/jki/pkg/factory"
	"github.com/spf13/cobra"
)

type Commander func(factory.Factory) *cobra.Command

func main() {
	cf := configflags.New()
	f := factory.New(cf)

	rootCmd := cobra.Command{
		Use: "jki",
	}
	cf.AddFlags(rootCmd.PersistentFlags())

	commanders := []Commander{
		build.NewCmdBuild,
		completion.NewCmdCompletion,
		config.NewCmdConfig,
		cp.NewCmdCp,
		transferimage.NewCmdTransferImage,
		version.NewCmdVersion,
	}

	for _, c := range commanders {
		rootCmd.AddCommand(c(f))
	}

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
