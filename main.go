package main

import (
	"log"

	"github.com/iftechio/jki/pkg/cmd/build"
	"github.com/iftechio/jki/pkg/cmd/completion"
	"github.com/iftechio/jki/pkg/cmd/config"
	"github.com/iftechio/jki/pkg/cmd/configflags"
	"github.com/iftechio/jki/pkg/cmd/cp"
	"github.com/iftechio/jki/pkg/cmd/transferimage"
	cmdutils "github.com/iftechio/jki/pkg/cmd/utils"
	"github.com/iftechio/jki/pkg/cmd/version"
	"github.com/spf13/cobra"
)

type Commander func(cmdutils.Factory) *cobra.Command

func main() {
	cf := configflags.New()
	factory := cmdutils.NewFactory(cf)

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
		rootCmd.AddCommand(c(factory))
	}

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
