package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/metafates/mangal/common"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/history"
	"github.com/metafates/mangal/style"
	"github.com/metafates/mangal/tui"
	"github.com/metafates/mangal/util"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

var mangalCmd = &cobra.Command{
	Use:   strings.ToLower(common.Mangal),
	Short: common.Mangal + " - A Manga Downloader",
	Long: common.AsciiArt + `

The ultimate CLI manga downloader`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("config")
		resume, _ := cmd.Flags().GetBool("resume")
		incognito, _ := cmd.Flags().GetBool("incognito")
		initConfig(configPath, false)

		if format, _ := cmd.Flags().GetString("format"); format != "" {
			config.UserConfig.Formats.Default = common.FormatType(format)
		}

		err := config.ValidateConfig(config.UserConfig)
		if err != nil {
			log.Fatal(err)
		}

		if !incognito {
			initAnilist()
		} else {
			config.UserConfig.IncognitoMode = true
			config.UserConfig.Anilist.Enabled = false
		}

		var (
			bubble  *tui.Bubble
			options []tea.ProgramOption
		)

		if config.UserConfig.UI.Fullscreen {
			options = append(options, tea.WithAltScreen())
		} else {
			style.CommonStyle.Margin(1, 1)
		}

		if resume {
			config.UserConfig.HistoryMode = true

			bubble = tui.NewBubble(tui.ResumeState)

			h, err := history.ReadHistory()

			if err != nil {
				log.Fatal(err)
			}

			var items []list.Item

			for _, item := range h {
				items = append(items, item)
			}

			sort.Slice(items, func(i, j int) bool {
				return items[i].(*history.Entry).Manga.Info < items[j].(*history.Entry).Manga.Info
			})

			bubble.ResumeList.SetItems(items)
		} else {
			bubble = tui.NewBubble(tui.SearchState)
		}

		program := tea.NewProgram(bubble, options...)

		if err := program.Start(); err != nil {
			log.Fatal(err)
		}
	},
}

// initConfig initializes the config file
// If the given string is empty, it will use the default config file
func initConfig(configPath string, validate bool) {
	if configPath != "" {
		// check if config is a TOML file
		if filepath.Ext(configPath) != ".toml" {
			log.Fatal("config file must be a TOML file")
		}

		// check if config file exists
		exists, err := afero.Exists(filesystem.Get(), configPath)

		if err != nil {
			log.Fatal(errors.New("access to config file denied"))
		}

		// if config file doesn't exist raise error
		configPath = path.Clean(configPath)
		if !exists {
			log.Fatal(errors.New(fmt.Sprintf("config at path %s doesn't exist", configPath)))
		}

		config.UserConfig = config.GetConfig(configPath)
	} else {
		// if config path is empty, use default config file
		config.UserConfig = config.GetConfig("")
	}

	if !validate {
		return
	}

	// check if config file is valid
	err := config.ValidateConfig(config.UserConfig)

	if err != nil {
		log.Fatal(err)
	}
}

func initAnilist() {
	if config.UserConfig == nil {
		log.Fatal("config is not initialized")
	}

	// check if anilist is enabled and token is expired
	if config.UserConfig.Anilist.Enabled && config.UserConfig.Anilist.Client.IsExpired() {
		fmt.Println("You are seeing this because you have enabled Anilist integration")
		fmt.Println()
		fmt.Printf("Anilist token is expired, press %s to open anilist page with a new token\n", style.AccentStyle.Render("enter"))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		fmt.Println("Opening Anilist page...")
		err := open.Run(config.UserConfig.Anilist.Client.AuthURL())

		if err != nil {
			fmt.Println("Something went wrong, please copy the url below manually")
			fmt.Println(style.AccentStyle.Render(config.UserConfig.Anilist.Client.AuthURL()))
		}

		// wait for user to input token
		fmt.Println()
		fmt.Print("Enter token: ")

		if scanner.Scan() {
			token := scanner.Text()
			if err := config.UserConfig.Anilist.Client.Login(token); err != nil {
				log.Fatal("could not login to Anilist. Are you using the correct token?")
			}
		} else {
			log.Fatal("could not read token")
		}
	}
}

func init() {
	mangalCmd.Flags().StringP("config", "c", "", "use config from path")
	mangalCmd.Flags().StringP("format", "f", "", "use custom format")
	mangalCmd.Flags().BoolP("incognito", "i", false, "do not save history")
	mangalCmd.Flags().BoolP("resume", "r", false, "resume reading")
	_ = mangalCmd.MarkFlagFilename("config", "toml")
	_ = mangalCmd.RegisterFlagCompletionFunc("format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return util.Map(common.AvailableFormats, util.ToString[common.FormatType]), cobra.ShellCompDirectiveDefault
	})
}

// Execute executes root command
func Execute() {
	_ = mangalCmd.Execute()
}
