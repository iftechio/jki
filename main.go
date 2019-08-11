package main

import (
	"log"

	"github.com/bario/jki/pkg/cmd/build"
	"github.com/bario/jki/pkg/cmd/config"
	"github.com/bario/jki/pkg/cmd/configflags"
	"github.com/bario/jki/pkg/cmd/cp"
	cmdutils "github.com/bario/jki/pkg/cmd/utils"
	"github.com/bario/jki/pkg/cmd/version"
	"github.com/spf13/cobra"
)

func main() {
	cf := configflags.New()
	factory := cmdutils.NewFactory(cf)

	rootCmd := cobra.Command{}
	cf.AddFlags(rootCmd.PersistentFlags())

	rootCmd.AddCommand(cp.NewCmdCp(factory))
	rootCmd.AddCommand(config.NewCmdConfig(factory))
	rootCmd.AddCommand(build.NewCmdBuild(factory))
	rootCmd.AddCommand(version.NewCmdVersion())

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
