package cmd

import (
	"bufio"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/metafates/mangal/common"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/history"
	"github.com/metafates/mangal/style"
	"github.com/metafates/mangal/tui"
	"github.com/metafates/mangal/util"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"log"
	"os"
	"sort"
	"strings"
)

var mangalCmd = &cobra.Command{
	Use:   strings.ToLower(common.Mangal),
	Short: common.Mangal + " - A Manga Downloader",
	Long: style.Bold.Render(common.AsciiArt) + `

The ultimate CLI manga downloader`,
	Run: func(cmd *cobra.Command, args []string) {
		resume, _ := cmd.Flags().GetBool("resume")
		incognito, _ := cmd.Flags().GetBool("incognito")

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
			style.Common.Margin(1, 1)
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

func initAnilist() {
	if config.UserConfig == nil {
		log.Fatal("config is not initialized")
	}

	// check if anilist is enabled and token is expired
	if config.UserConfig.Anilist.Enabled && config.UserConfig.Anilist.Client.IsExpired() {
		fmt.Println("You are seeing this because you have enabled Anilist integration")
		fmt.Println()
		fmt.Printf("Anilist token is expired, press %s to open anilist page with a new token\n", style.Accent.Render("enter"))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		fmt.Println("Opening Anilist page...")
		err := open.Run(config.UserConfig.Anilist.Client.AuthURL())

		if err != nil {
			fmt.Println("Something went wrong, please copy the url below manually")
			fmt.Println(style.Accent.Render(config.UserConfig.Anilist.Client.AuthURL()))
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
	mangalCmd.Flags().StringP("format", "f", "", "use custom format")
	mangalCmd.Flags().BoolP("incognito", "i", false, "do not save history")
	mangalCmd.Flags().BoolP("resume", "r", false, "resume reading")
	_ = mangalCmd.RegisterFlagCompletionFunc("format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return util.Map(common.AvailableFormats, util.ToString[common.FormatType]), cobra.ShellCompDirectiveDefault
	})

	// colorize help message
	cobra.AddTemplateFunc("StyleHeading", lipgloss.NewStyle().Foreground(lipgloss.Color("#FAE3B0")).Render)
	usageTemplate := mangalCmd.UsageTemplate()
	usageTemplate = strings.NewReplacer(
		`Usage:`, `{{StyleHeading "Usage:"}}`,
		`Aliases:`, `{{StyleHeading "Aliases:"}}`,
		`Available Commands:`, `{{StyleHeading "Available Commands:"}}`,
		`Global Flags:`, `{{StyleHeading "Global Flags:"}}`,
		`Flags:`, `{{StyleHeading "Flags:"}}`,
	).Replace(usageTemplate)
	mangalCmd.SetUsageTemplate(usageTemplate)
}

// Execute executes root command
func Execute() {
	// init config
	config.Initialize("", true)

	_ = mangalCmd.Execute()
}
