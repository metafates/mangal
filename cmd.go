package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   strings.ToLower(Mangal),
	Short: Mangal + " - A Manga Downloader",
	Long: AsciiArt + `

The ultimate CLI manga downloader`,
	Run: func(cmd *cobra.Command, args []string) {
		config, _ := cmd.Flags().GetString("config")
		initConfig(config)

		if format, _ := cmd.Flags().GetString("format"); format != "" {
			UserConfig.Format = FormatType(format)
		}

		if err := ValidateConfig(UserConfig); err != nil {
			log.Fatal(err)
		}

		var program *tea.Program

		if UserConfig.UI.Fullscreen {
			program = tea.NewProgram(NewBubble(searchState), tea.WithAltScreen())
		} else {
			program = tea.NewProgram(NewBubble(searchState))
		}

		if err := program.Start(); err != nil {
			log.Fatal(err)
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Long:  fmt.Sprintf("Shows %s versions and build date", Mangal),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s version %s\n", accentStyle.Render(strings.ToLower(Mangal)), boldStyle.Render(version))
	},
}

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Remove cached and temp files",
	Long:  "Removes cached files produced by scraper and temp files from downloader",
	Run: func(cmd *cobra.Command, args []string) {
		counter, bytes := RemoveTemp()
		c, b := RemoveCache()
		counter += c
		bytes += b

		fmt.Printf("%d files removed\nCleaned up %.2fMB\n", counter, BytesToMegabytes(bytes))
	},
}

var cleanupTempCmd = &cobra.Command{
	Use:   "temp",
	Short: "Remove temp files",
	Long:  "Removes temp files produced by downloader",
	Run: func(cmd *cobra.Command, args []string) {
		counter, bytes := RemoveTemp()
		fmt.Printf("%d temp files removed\nCleaned up %.2fMB\n", counter, BytesToMegabytes(bytes))
	},
}

var cleanupCacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Remove cache files",
	Long:  "Removes cache files produced by scraper",
	Run: func(cmd *cobra.Command, args []string) {
		counter, bytes := RemoveCache()
		fmt.Printf("%d cache files removed\nCleaned up %.2fMB\n", counter, BytesToMegabytes(bytes))
	},
}

var inlineCmd = &cobra.Command{
	Use:   "inline",
	Short: "Search & Download manga in inline mode",
	Long: `Search & Download manga in inline mode
Useful for scripting`,
	Run: func(cmd *cobra.Command, args []string) {
		query, _ := cmd.Flags().GetString("query")
		mangaIdx, _ := cmd.Flags().GetInt("manga")
		chapterIdx, _ := cmd.Flags().GetInt("chapter")
		asJson, _ := cmd.Flags().GetBool("json")
		format, _ := cmd.Flags().GetString("format")
		showUrls, _ := cmd.Flags().GetBool("urls")
		asTemp, _ := cmd.Flags().GetBool("temp")
		doRead, _ := cmd.Flags().GetBool("read")
		doOpen, _ := cmd.Flags().GetBool("open")
		config, _ := cmd.Flags().GetString("config")

		res, err := InlineMode(query, InlineOptions{
			config:     config,
			mangaIdx:   mangaIdx,
			chapterIdx: chapterIdx,
			asJson:     asJson,
			format:     FormatType(format),
			showUrls:   showUrls,
			asTemp:     asTemp,
			doRead:     doRead,
			doOpen:     doOpen,
		})

		if err != nil {
			if asJson {
				fmt.Printf(`{error: "%s"}\n`, err)
			} else {
				fmt.Println(err)
			}

			os.Exit(1)
		}

		fmt.Print(res)
	},
}

// initConfig initializes the config file
// If the given string is empty, it will use the default config file
func initConfig(config string) {
	exists, err := Afero.Exists(config)

	if err != nil {
		log.Fatal(errors.New("access to config file denied"))
	}

	if config != "" {
		config = path.Clean(config)
		if !exists {
			log.Fatal(errors.New(fmt.Sprintf("config at path %s doesn't exist", config)))
		}

		UserConfig = GetConfig(config)
	} else {
		UserConfig = GetConfig("") // get config from default config path
	}
}

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
		configPath, err := GetConfigPath()

		if err != nil {
			log.Fatal("Can't get config location, permission denied, probably")
		}

		exists, err := Afero.Exists(configPath)

		if err != nil {
			log.Fatalf("Can't understand if config exists or not. It is expected at\n%s\n", configPath)
		}

		if exists {
			fmt.Printf("Config exists at\n%s\n", successStyle.Render(configPath))
		} else {
			fmt.Printf("Config doesn't exist, but it is expected to be at\n%s\n", successStyle.Render(configPath))
		}
	},
}

var configPreviewCmd = &cobra.Command{
	Use:   "preview",
	Short: "Preview current config",
	Run: func(cmd *cobra.Command, args []string) {
		configPath, err := GetConfigPath()

		if err != nil {
			log.Fatal("Permission to config file was denied")
		}

		exists, err := Afero.Exists(configPath)

		if err != nil {
			log.Fatalf("Permission to config file was denied")
		}

		if !exists {
			log.Fatal("Config doesn't exist")
		}

		contents, err := Afero.ReadFile(configPath)
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
		configPath, err := GetConfigPath()

		if err != nil {
			log.Fatal("Permission to config file was denied")
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

		configPath, err := GetConfigPath()

		if err != nil {
			log.Fatal("Can't get config location, permission denied, probably")
		}

		exists, err := Afero.Exists(configPath)

		var createConfig = func() {
			if err := Afero.MkdirAll(filepath.Dir(configPath), 0700); err != nil {
				log.Fatal("Error while creating file")
			} else if file, err := Afero.Create(configPath); err != nil {
				log.Fatal("Error while creating file")
			} else if _, err = file.Write(DefaultConfigBytes); err != nil {
				log.Fatal("Error while writing to file")
			} else {
				fmt.Printf("Config created at\n%s\n", successStyle.Render(configPath))
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
			log.Fatalf("Config file already exists. Use %s to overwrite it", accentStyle.Render("--force"))
		} else {
			createConfig()
		}
	},
}

var configTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test user config for any errors",
	Run: func(cmd *cobra.Command, args []string) {
		config, _ := cmd.Flags().GetString("config")
		initConfig(config)

		if err := ValidateConfig(UserConfig); err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(successStyle.Render("Everything is OK"))
		}
	},
}

var formatsCmd = &cobra.Command{
	Use:   "formats",
	Short: "Information about available formats",
	Long:  "Show information about available formats with quick description of each",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(boldStyle.Render("Available formats") + "\n\n")
		for _, format := range AvailableFormats {
			fmt.Printf("%s - %s\n", accentStyle.Render(string(format)), FormatsInfo[format])
		}
	},
}

var checkUpdateCmd = &cobra.Command{
	Use:   "check-update",
	Short: "Check if new version is available",
	Long:  "Fethces latest version of the program from github and compares it with current version",
	Run: func(cmd *cobra.Command, args []string) {
		const githubReleaseURL = "https://github.com/metafates/mangal/releases"

		resp, err := http.Get(githubReleaseURL)
		if err != nil {
			log.Fatal(err)
		}

		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		var latestVersion string

		// compile regex for tag
		re := regexp.MustCompile(`^v\d+\.\d+\.\d+$`)

		// get latest release tag
		doc.Find("a > div > span").Each(func(_ int, s *goquery.Selection) {
			var tag = strings.TrimSpace(s.Text())

			// check if latestVersion matches tag regex and set it if it does (only if it is not set yet)
			if re.MatchString(tag) && latestVersion == "" {
				// remove "v" from tag
				latestVersion = strings.TrimPrefix(tag, "v")
			}
		})

		if latestVersion == "" {
			log.Fatalf("Can't find latest version\nYou can visit %s to check for updates", githubReleaseURL)
		}

		// check if current version is latest
		if latestVersion == version {
			fmt.Printf("You are using the latest version of %s\n", Mangal)
		} else {
			fmt.Printf("New version of %s is available: %s\n", Mangal, accentStyle.Render(latestVersion))
			fmt.Printf("You can download it from %s\n", accentStyle.Render(githubReleaseURL))
			fmt.Println("Or use your package manager to update")
		}
	},
}

// CmdExecute adds all subcommands to the root command and executes it
func CmdExecute() {
	rootCmd.AddCommand(versionCmd)

	cleanupCmd.AddCommand(cleanupTempCmd)
	cleanupCmd.AddCommand(cleanupCacheCmd)
	rootCmd.AddCommand(cleanupCmd)

	configCmd.AddCommand(configWhereCmd)
	configCmd.AddCommand(configPreviewCmd)
	configCmd.AddCommand(configEditCmd)
	configInitCmd.Flags().BoolP("force", "f", false, "overwrite existing config")
	configCmd.AddCommand(configInitCmd)

	configTestCmd.Flags().StringP("config", "c", "", "use config from path")
	configCmd.AddCommand(configTestCmd)
	rootCmd.AddCommand(configCmd)

	inlineCmd.Flags().Int("manga", -1, "choose manga by index")
	inlineCmd.Flags().Int("chapter", -1, "choose and download chapter by index")
	inlineCmd.Flags().StringP("query", "q", "", "manga to search")
	inlineCmd.Flags().BoolP("json", "j", false, "print as json")
	inlineCmd.Flags().StringP("format", "f", "", "use custom format")
	inlineCmd.Flags().BoolP("urls", "u", false, "show urls")
	inlineCmd.Flags().BoolP("temp", "t", false, "download as temp")
	inlineCmd.Flags().BoolP("read", "r", false, "read chapter")
	inlineCmd.Flags().BoolP("open", "o", false, "open url")
	inlineCmd.Flags().StringP("config", "c", "", "use config from path")
	inlineCmd.Flags().SortFlags = false
	rootCmd.AddCommand(inlineCmd)

	rootCmd.Flags().StringP("config", "c", "", "use config from path")
	rootCmd.Flags().StringP("format", "f", "", "use custom format")

	rootCmd.AddCommand(formatsCmd)
	rootCmd.AddCommand(checkUpdateCmd)

	_ = rootCmd.Execute()
}
