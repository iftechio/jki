package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/iftechio/jki/pkg/cmd/build"
	"github.com/iftechio/jki/pkg/cmd/completion"
	"github.com/iftechio/jki/pkg/cmd/config"
	"github.com/iftechio/jki/pkg/cmd/cp"
	"github.com/iftechio/jki/pkg/cmd/transferimage"
	"github.com/iftechio/jki/pkg/cmd/upgrade"
	"github.com/iftechio/jki/pkg/cmd/version"
	"github.com/iftechio/jki/pkg/configflags"
	"github.com/iftechio/jki/pkg/factory"
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
		upgrade.NewCmdUpgrade,
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
