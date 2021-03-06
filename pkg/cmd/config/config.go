package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/iftechio/jki/pkg/factory"
	"github.com/iftechio/jki/pkg/registry"
)

func NewCmdConfig(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Modify config file",
	}

	saveDefaultConfig := false
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Default settings",
		Run: func(cmd *cobra.Command, args []string) {
			if !saveDefaultConfig {
				fmt.Print(defaultConfig)
				return
			}
			configPath := f.ConfigPath()
			_, err := os.Stat(configPath)
			if err == nil {
				fmt.Printf("%s already exists\n", configPath)
				return
			}
			err = ioutil.WriteFile(configPath, []byte(defaultConfig), 0644)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	initCmd.Flags().BoolVar(&saveDefaultConfig, "save", false, "Save default settings")

	editCmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit the config file from the default editor",
		RunE: func(cmd *cobra.Command, args []string) error {
			editor := os.Getenv("EDITOR")
			if len(editor) == 0 {
				switch runtime.GOOS {
				case "darwin", "linux":
					editor = "vi"
				case "windows":
					editor = "notepad"
				}
			}
			configPath := f.ConfigPath()
			c := exec.Command(editor, configPath)
			c.Stdin = os.Stdin
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			return c.Run()
		},
	}

	viewCmd := &cobra.Command{
		Use:   "view",
		Short: "Display settings from the config file",
		Run: func(cmd *cobra.Command, args []string) {
			configPath := f.ConfigPath()
			data, err := ioutil.ReadFile(configPath)
			if err != nil {
				log.Fatal(err)
			}
			_, _ = os.Stdout.Write(data)
		},
	}

	checkCmd := &cobra.Command{
		Use:   "check",
		Short: "Test configuration and exit",
		Run: func(cmd *cobra.Command, args []string) {
			configPath := f.ConfigPath()
			_, regs, err := registry.LoadRegistries(configPath)
			if err != nil {
				log.Fatal(err)
			}
			for name, reg := range regs {
				if err := reg.Verify(); err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "%s: invalid config: %s\n", name, err)
					return
				}
			}
			fmt.Println("OK!")
		},
	}

	getRegsCmd := &cobra.Command{
		Use:   "get-registries",
		Short: "List all registries defined in the config",
		Run: func(cmd *cobra.Command, args []string) {
			configPath := f.ConfigPath()
			_, regs, err := registry.LoadRegistries(configPath)
			if err != nil {
				log.Fatal(err)
			}
			for name := range regs {
				fmt.Println(name)
			}
		},
	}

	cmd.AddCommand(initCmd)
	cmd.AddCommand(editCmd)
	cmd.AddCommand(viewCmd)
	cmd.AddCommand(checkCmd)
	cmd.AddCommand(getRegsCmd)
	return cmd
}
