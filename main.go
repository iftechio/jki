package main

import (
	"log"

	"github.com/bario/jki/pkg/cmd/build"
	"github.com/bario/jki/pkg/cmd/completion"
	"github.com/bario/jki/pkg/cmd/config"
	"github.com/bario/jki/pkg/cmd/cp"
	"github.com/bario/jki/pkg/cmd/transferimage"
	"github.com/bario/jki/pkg/cmd/version"
	"github.com/bario/jki/pkg/configflags"
	"github.com/bario/jki/pkg/factory"
	"github.com/spf13/cobra"
)

type Commander func(factory.Factory) *cobra.Command

func main() {
	cf := configflags.New()
	factory := factory.New(cf)

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
