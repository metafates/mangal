package cmd

import (
	"fmt"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/style"
	"github.com/metafates/mangal/util"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Config actions",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var configWhereCmd = &cobra.Command{
	Use:   "where",
	Short: "Show config location",
	Long:  "Show path where config is located if it exists.\nOtherwise show path where it is expected to be",
	Run: func(cmd *cobra.Command, args []string) {
		configPath, err := util.UserConfigFile()

		if err != nil {
			log.Fatal("Can't get config location, permission denied, probably")
		}

		exists, err := afero.Exists(filesystem.Get(), configPath)

		if err != nil {
			log.Fatalf("Can't understand if config exists or not. It is expected at\n%s\n", configPath)
		}

		if exists {
			fmt.Printf("Config exists at\n%s\n", style.Success.Render(configPath))
		} else {
			fmt.Printf("Config doesn't exist, but it is expected to be at\n%s\n", style.Success.Render(configPath))
		}
	},
}

var configPreviewCmd = &cobra.Command{
	Use:   "preview",
	Short: "Preview current config",
	Long:  "Preview current config.\nIt will use `bat` to preview the config file if possible",
	Run: func(cmd *cobra.Command, args []string) {
		configPath, err := util.UserConfigFile()

		if err != nil {
			log.Fatal("Permission to config file was denied")
		}

		exists, err := afero.Exists(filesystem.Get(), configPath)

		if err != nil {
			log.Fatal("Permission to config file was denied")
		}

		if !exists {
			log.Fatal("Config doesn't exist")
		}

		// check if bat command is installed
		_, err = exec.LookPath("bat")
		if err == nil {
			cmd := exec.Command("bat", "-l", "toml", configPath)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			return
		}

		// check if less command is installed
		_, err = exec.LookPath("less")
		if err == nil {
			cmd := exec.Command("less", configPath)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			return
		}

		contents, err := afero.ReadFile(filesystem.Get(), configPath)
		if err != nil {
			log.Fatal("Permission to config file was denied")
		}

		fmt.Println(string(contents))
	},
}

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit config in the default editor",
	Long:  "Edit config in the default editor.\nIf config doesn't exist, it will be created",
	Run: func(cmd *cobra.Command, args []string) {
		configPath, err := util.UserConfigFile()

		if err != nil {
			log.Fatal("Permission to config file was denied")
		}

		// check if config file exists
		exists, err := afero.Exists(filesystem.Get(), configPath)
		if err != nil {
			log.Fatal("Permission to config file was denied")
		}

		if !exists {
			fmt.Println("Config doesn't exist, nothing to edit")
			os.Exit(0)
		}

		err = open.Start(configPath)
		if err != nil {
			log.Fatal("Can't open editor")
		}
	},
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Init default config",
	Long:  "Init default config at the default location.\nIf the config already exists, it will not be overwritten",
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")
		clean, _ := cmd.Flags().GetBool("clean")

		configPath, err := util.UserConfigFile()

		if err != nil {
			log.Fatal("Can't get config location, permission denied, probably")
		}

		exists, err := afero.Exists(filesystem.Get(), configPath)

		var createConfig = func() {
			var configToWrite string

			if clean {
				// remove all lines with comments from toml string
				configToWrite = regexp.MustCompile("\n[^\n]*#.*").ReplaceAllString(string(config.DefaultConfigBytes), "")

				// remove all empty lines from toml string
				configToWrite = regexp.MustCompile("\n\n+").ReplaceAllString(configToWrite, "\n")

				// insert newline before each section
				configToWrite = regexp.MustCompile("(?m)^(\\[.*])").ReplaceAllString(configToWrite, "\n$1")
			} else {
				configToWrite = string(config.DefaultConfigBytes)
			}

			if err := filesystem.Get().MkdirAll(filepath.Dir(configPath), 0700); err != nil {
				log.Fatal("Error while creating file")
			} else if file, err := filesystem.Get().Create(configPath); err != nil {
				log.Fatal("Error while creating file")
			} else if _, err = file.Write([]byte(configToWrite)); err != nil {
				log.Fatal("Error while writing to file")
			} else {
				fmt.Printf("Config created at\n%s\n", style.Success.Render(configPath))
			}
		}

		if force {
			createConfig()
			return
		}

		if err != nil {
			log.Fatalf("Can't understand if config exists or not, but it is expected at\n%s\n", configPath)
		}

		if exists {
			log.Fatalf("Config file already exists. Use %s to overwrite it", style.Accent.Render("--force"))
		} else {
			createConfig()
		}
	},
}

var configRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove config",
	Long:  "Remove config.\nIf config doesn't exist, it will not be removed",
	Run: func(cmd *cobra.Command, args []string) {
		configPath, err := util.UserConfigFile()

		if err != nil {
			log.Fatal("Can't get config location, permission denied, probably")
		}

		exists, err := afero.Exists(filesystem.Get(), configPath)

		if err != nil {
			log.Fatalf("Can't understand if config exists or not. It is expected at\n%s\n", configPath)
		}

		if exists {
			if err := filesystem.Get().Remove(configPath); err != nil {
				log.Fatal("Error while removing file")
			} else {
				fmt.Println("Config removed")
			}
		} else {
			fmt.Println("Config doesn't exist, nothing to remove")
		}
	},
}

func init() {
	configCmd.AddCommand(configWhereCmd)
	configCmd.AddCommand(configPreviewCmd)
	configCmd.AddCommand(configEditCmd)
	configCmd.AddCommand(configInitCmd)

	configInitCmd.Flags().BoolP("force", "f", false, "overwrite existing config")
	configInitCmd.Flags().BoolP("clean", "c", false, "do not add comments and empty lines")

	configCmd.AddCommand(configRemoveCmd)
	mangalCmd.AddCommand(configCmd)
}
