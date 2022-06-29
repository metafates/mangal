package main

import (
	"bufio"
	"errors"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"os/exec"
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
		incognito, _ := cmd.Flags().GetBool("incognito")
		initConfig(config, false)

		if !incognito {
			initAnilist()
		} else {
			UserConfig.Anilist.Enabled = false
		}

		if format, _ := cmd.Flags().GetString("format"); format != "" {
			UserConfig.Formats.Default = FormatType(format)
		}

		err := ValidateConfig(UserConfig)
		if err != nil {
			log.Fatal(err)
		}

		var program *tea.Program

		if UserConfig.UI.Fullscreen {
			program = tea.NewProgram(NewBubble(searchState), tea.WithAltScreen())
		} else {
			commonStyle.Margin(1, 1)
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
		fmt.Printf("%s version %s\n", Mangal, accentStyle.Render(Version))
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

var cleanupAnilistCmd = &cobra.Command{
	Use:   "anilist",
	Short: "Remove Anilist cache",
	Run: func(cmd *cobra.Command, args []string) {
		// get config dir
		configDir, err := os.UserConfigDir()

		if err != nil {
			log.Fatal(err)
		}

		// get anilist file
		anilistFile := filepath.Join(configDir, Mangal, "anilist.json")

		// check if anilist file exists
		exists, err := Afero.Exists(anilistFile)

		if err != nil {
			log.Fatal(err)
		}

		// if anilist file doesn't exist exit
		if !exists {
			fmt.Println("Anilist file doesn't exist so nothing to clean up")
			return
		}

		// remove anilist file
		err = Afero.Remove(anilistFile)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Anilist file removed")
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
func initConfig(config string, validate bool) {
	if config != "" {
		// check if config is a TOML file
		if filepath.Ext(config) != ".toml" {
			log.Fatal("config file must be a TOML file")
		}

		// check if config file exists
		exists, err := Afero.Exists(config)

		if err != nil {
			log.Fatal(errors.New("access to config file denied"))
		}

		// if config file doesn't exist raise error
		config = path.Clean(config)
		if !exists {
			log.Fatal(errors.New(fmt.Sprintf("config at path %s doesn't exist", config)))
		}

		UserConfig = GetConfig(config)
	} else {
		// if config path is empty, use default config file
		UserConfig = GetConfig("")
	}

	if !validate {
		return
	}

	// check if config file is valid
	err := ValidateConfig(UserConfig)

	if err != nil {
		log.Fatal(err)
	}
}

func initAnilist() {
	if UserConfig == nil {
		log.Fatal("config is not initialized")
	}

	// check if anilist is enabled and token is expired
	if UserConfig.Anilist.Enabled && UserConfig.Anilist.Client.IsExpired() {
		fmt.Println("You are seeing this because you have enabled Anilist integration")
		fmt.Println()
		fmt.Printf("Anilist token is expired, press %s to open anilist page with a new token\n", accentStyle.Render("enter"))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		fmt.Println("Opening Anilist page...")
		err := open.Run(UserConfig.Anilist.Client.AuthURL())

		if err != nil {
			fmt.Println("Something went wrong, please copy the url below manually")
			fmt.Println(accentStyle.Render(UserConfig.Anilist.Client.AuthURL()))
		}

		// wait for user to input token
		fmt.Println()
		fmt.Print("Enter token: ")

		if scanner.Scan() {
			token := scanner.Text()
			if err := UserConfig.Anilist.Client.Login(token); err != nil {
				log.Fatal("could not login to Anilist. Are you using the correct token?")
			}
		} else {
			log.Fatal("could not read token")
		}
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
	Long:  "Preview current config.\nIt will use `bat` to preview the config file if possible",
	Run: func(cmd *cobra.Command, args []string) {
		configPath, err := GetConfigPath()

		if err != nil {
			log.Fatal("Permission to config file was denied")
		}

		exists, err := Afero.Exists(configPath)

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
		clean, _ := cmd.Flags().GetBool("clean")

		configPath, err := GetConfigPath()

		if err != nil {
			log.Fatal("Can't get config location, permission denied, probably")
		}

		exists, err := Afero.Exists(configPath)

		var createConfig = func() {
			var configToWrite string

			if clean {
				// remove all lines with comments from toml string
				configToWrite = regexp.MustCompile("\n[^\n]*#.*").ReplaceAllString(DefaultConfigString, "")

				// remove all empty lines from toml string
				configToWrite = regexp.MustCompile("\n\n+").ReplaceAllString(configToWrite, "\n")

				// insert newline before each section
				configToWrite = regexp.MustCompile("(?m)^(\\[.*])").ReplaceAllString(configToWrite, "\n$1")
			} else {
				configToWrite = DefaultConfigString
			}

			if err := Afero.MkdirAll(filepath.Dir(configPath), 0700); err != nil {
				log.Fatal("Error while creating file")
			} else if file, err := Afero.Create(configPath); err != nil {
				log.Fatal("Error while creating file")
			} else if _, err = file.Write([]byte(configToWrite)); err != nil {
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

var configRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove config",
	Long:  "Remove config.\nIf config doesn't exist, it will not be removed",
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
			if err := Afero.Remove(configPath); err != nil {
				log.Fatal("Error while removing file")
			} else {
				fmt.Println("Config removed")
			}
		} else {
			fmt.Println("Config doesn't exist, nothing to remove")
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

var latestCmd = &cobra.Command{
	Use:   "latest",
	Short: fmt.Sprintf("Check if latest version of %s is used", Mangal),
	Long:  "Fetches the latest version from the GitHub and compares it with current version",
	Run: func(cmd *cobra.Command, args []string) {
		const githubReleaseURL = "https://github.com/metafates/mangal/releases/latest"

		latestVersion, err := FetchLatestVersion()

		if err != nil || latestVersion == "" {
			log.Fatalf("Can't find latest version\nYou can visit %s to check for updates", githubReleaseURL)
		}

		// check if current version is latest
		if latestVersion <= Version {
			fmt.Printf("You are using the latest version of %s\n", Mangal)
		} else {
			fmt.Printf("New version of %s is available: %s\n", Mangal, accentStyle.Render(latestVersion))
			fmt.Printf("You can download it from %s\n", accentStyle.Render(githubReleaseURL))
			fmt.Println("Or use your package manager to update")
		}
	},
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Run this in case of any errors",
	Long: `Check if ` + Mangal + ` is properly configured.
It checks if config file is valid and used sources are available`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			ok = func() {
				fmt.Print(successStyle.Render("OK"))
				fmt.Println()
			}

			fail = func() {
				fmt.Print(failStyle.Render("Fail"))
				fmt.Println()
			}
		)

		fmt.Print("Checking if latest version is used... ")
		latestVersion, err := FetchLatestVersion()
		if err != nil {
			fail()
			fmt.Printf("Can't find latest version\nRun %s to get more information\n", accentStyle.Render("mangal latest"))
			os.Exit(1)
		} else if latestVersion > Version {
			fail()
			fmt.Printf("New version of %s is available: %s\n", Mangal, accentStyle.Render(latestVersion))
			fmt.Printf("Run %s to get more information\n", accentStyle.Render("mangal latest"))
			os.Exit(1)
		} else {
			ok()
		}

		fmt.Print("Checking config... ")
		UserConfig = GetConfig("")

		err = ValidateConfig(UserConfig)
		if err != nil {
			fail()
			fmt.Printf("Config error: %s\n", err)
			os.Exit(1)
		}

		ok()

		var sourceNotAvailable = func(source *Source) {
			fail()
			fmt.Printf("Source %s is not available\n", source.Name)
			fmt.Printf("Try to reinitialize your config with %s\n", accentStyle.Render("mangal config init --force"))
			fmt.Println("Note, that this will overwrite your current config")
			os.Exit(1)
		}

		scanner := bufio.NewScanner(os.Stdin)

		// check if scraper sources are online
		for _, scraper := range UserConfig.Scrapers {
			source := scraper.Source

			// read line from stdin
			fmt.Printf("Please, enter a manga title to test %s: ", source.Name)
			scanner.Scan()

			if scanner.Err() != nil {
				fail()
				fmt.Printf("Error while reading from stdin: %s\n", err)
				os.Exit(1)
			}

			fmt.Printf("Checking source %s... ", source.Name)
			resp, err := http.Get(source.Base)
			if err != nil {
				sourceNotAvailable(source)
			}

			_ = resp.Body.Close()

			// check if response is 200
			if resp.StatusCode != 200 {
				sourceNotAvailable(source)
			}

			// try to get any manga page using scraper
			manga, err := scraper.SearchManga(scanner.Text())
			if err != nil {
				sourceNotAvailable(source)
			}

			// check if manga is not empty
			if len(manga) == 0 {
				sourceNotAvailable(source)
			}

			// get chapters for first manga
			chapters, err := scraper.GetChapters(manga[0])
			if err != nil {
				sourceNotAvailable(source)
			}

			// check if chapters is not empty
			if len(chapters) == 0 {
				sourceNotAvailable(source)
			}

			// get pages for first chapter
			pages, err := scraper.GetPages(chapters[0])
			if err != nil {
				sourceNotAvailable(source)
			}

			// check if pages is not empty
			if len(pages) == 0 {
				sourceNotAvailable(source)
			}

			// try to download first page
			image, err := scraper.GetFile(pages[0])
			if err != nil {
				sourceNotAvailable(source)
			}

			// check if images is not empty
			if image.Len() == 0 {
				sourceNotAvailable(source)
			}

			ok()
		}
	},
}

// Adds all child commands to the root command and sets flags appropriately.
func init() {
	rootCmd.AddCommand(versionCmd)

	cleanupCmd.AddCommand(cleanupTempCmd)
	cleanupCmd.AddCommand(cleanupCacheCmd)
	cleanupCmd.AddCommand(cleanupAnilistCmd)
	rootCmd.AddCommand(cleanupCmd)

	configCmd.AddCommand(configWhereCmd)
	configCmd.AddCommand(configPreviewCmd)
	configCmd.AddCommand(configEditCmd)
	configCmd.AddCommand(configInitCmd)

	configInitCmd.Flags().BoolP("force", "f", false, "overwrite existing config")
	configInitCmd.Flags().BoolP("clean", "c", false, "do not add comments and empty lines")

	configCmd.AddCommand(configRemoveCmd)
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
	_ = inlineCmd.MarkFlagRequired("query")
	_ = inlineCmd.MarkFlagFilename("config", "toml")
	_ = inlineCmd.RegisterFlagCompletionFunc("format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return Map(AvailableFormats, ToString[FormatType]), cobra.ShellCompDirectiveDefault
	})
	rootCmd.AddCommand(inlineCmd)

	rootCmd.Flags().StringP("config", "c", "", "use config from path")
	rootCmd.Flags().StringP("format", "f", "", "use custom format")
	rootCmd.Flags().BoolP("incognito", "i", false, "will not sync with anilist even if enabled")
	_ = rootCmd.MarkFlagFilename("config", "toml")
	_ = rootCmd.RegisterFlagCompletionFunc("format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return Map(AvailableFormats, ToString[FormatType]), cobra.ShellCompDirectiveDefault
	})

	rootCmd.AddCommand(formatsCmd)
	rootCmd.AddCommand(latestCmd)
	rootCmd.AddCommand(doctorCmd)
}

// CmdExecute executes root command
func CmdExecute() {
	_ = rootCmd.Execute()
}
