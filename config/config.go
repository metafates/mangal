package config

import (
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/where"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var EnvKeyReplacer = strings.NewReplacer(".", "_")

// Setup initializes the configuration
func Setup() error {
	setName()
	setFs()
	setEnvs()
	setDefaults()
	setPaths()

	err := viper.ReadInConfig()

	if err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			// Use defaults then
			return nil
		default:
			return err
		}
	}

	resolveAliases()
	return nil
}

func setName() {
	viper.SetConfigName(constant.Mangal)
	viper.SetConfigType("toml")
}

func setFs() {
	viper.SetFs(filesystem.Get())
}

// setPaths sets the paths to the config files
func setPaths() {
	viper.AddConfigPath(where.Config())
}

// setEnvs sets the environment variables
func setEnvs() {
	viper.SetEnvPrefix(constant.Mangal)
	viper.SetEnvKeyReplacer(EnvKeyReplacer)

	for _, env := range EnvExposed {
		viper.MustBindEnv(env)
	}
}

// setDefaults sets the default values
func setDefaults() {
	viper.SetTypeByDefaultValue(true)

	fields := map[string]any{
		// Downloader
		constant.DownloaderPath:                ".",
		constant.DownloaderChapterNameTemplate: "[{padded-index}] {chapter}",
		constant.DownloaderAsync:               true,
		constant.DownloaderCreateMangaDir:      true,
		constant.DownloaderDefaultSource:       "",
		constant.DownloaderStopOnError:         false,

		// Formats
		constant.FormatsUse:                   "pdf",
		constant.FormatsSkipUnsupportedImages: true,

		// Mini-mode
		constant.MiniSearchLimit: 20,

		// Icons
		constant.IconsVariant: "plain",

		// Reader
		constant.ReaderPDF:           "",
		constant.ReaderCBZ:           "",
		constant.ReaderZIP:           "",
		constant.RaderPlain:          "",
		constant.ReaderReadInBrowser: false,

		// History
		constant.HistorySaveOnRead:     true,
		constant.HistorySaveOnDownload: false,

		// Mangadex
		constant.MangadexLanguage:                "en",
		constant.MangadexNSFW:                    false,
		constant.MangadexShowUnavailableChapters: false,

		// Installer
		constant.InstallerUser:   "metafates",
		constant.InstallerRepo:   "mangal-scrapers",
		constant.InstallerBranch: "main",

		// Gen
		constant.GenAuthor: "",

		// Logs
		constant.LogsWrite: false,
		constant.LogsLevel: "info",

		// Anilist
		constant.AnilistEnable: false,
	}

	for field, value := range fields {
		viper.SetDefault(field, value)
	}
}

// resolveAliases resolves the aliases for the paths
func resolveAliases() {
	home := lo.Must(os.UserHomeDir())
	path := viper.GetString(constant.DownloaderPath)

	path = strings.ReplaceAll(path, "$HOME", home)
	path = strings.ReplaceAll(path, "~", home)

	viper.Set(constant.DownloaderPath, path)
}
