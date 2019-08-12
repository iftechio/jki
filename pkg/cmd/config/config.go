package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"

	cmdutils "github.com/iftechio/jki/pkg/cmd/utils"
	"github.com/iftechio/jki/pkg/registry"
)

func NewCmdConfig(f cmdutils.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Modify config file",
	}

	configPathFlag := "jkiconfig"
	saveDefaultConfig := false
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Default settings",
		Run: func(cmd *cobra.Command, args []string) {
			if !saveDefaultConfig {
				fmt.Print(defaultConfig)
				return
			}
			configPath := cmd.Flag(configPathFlag).Value.String()
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
				case "darwin":
					fallthrough
				case "linux":
					editor = "vi"
				case "windows":
					editor = "notepad"
				}
			}
			configPath := cmd.Flag(configPathFlag).Value.String()
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
			configPath := cmd.Flag(configPathFlag).Value.String()
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
			configPath := cmd.Flag(configPathFlag).Value.String()
			_, _, err := registry.LoadRegistries(configPath)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("OK!")
		},
	}

	cmd.AddCommand(initCmd)
	cmd.AddCommand(editCmd)
	cmd.AddCommand(viewCmd)
	cmd.AddCommand(checkCmd)
	return cmd
}
