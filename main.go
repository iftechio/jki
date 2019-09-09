package main

import (
	"log"

	"github.com/bario/jki/pkg/cmd/build"
	"github.com/bario/jki/pkg/cmd/completion"
	"github.com/bario/jki/pkg/cmd/config"
	"github.com/bario/jki/pkg/cmd/configflags"
	"github.com/bario/jki/pkg/cmd/cp"
	"github.com/bario/jki/pkg/cmd/transferimage"
	cmdutils "github.com/bario/jki/pkg/cmd/utils"
	"github.com/bario/jki/pkg/cmd/version"
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
