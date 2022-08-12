package config

import (
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/where"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"os"
	"runtime"
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

	switch err.(type) {
	case viper.ConfigFileNotFoundError:
		// Use defaults then
		return nil
	default:
		resolveAliases()
		return err
	}
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
		DownloaderPath:                ".",
		DownloaderChapterNameTemplate: "[{padded-index}] {chapter}",
		DownloaderAsync:               true,

		// Formats
		FormatsUse: "pdf",

		// Mini-mode
		MiniVimMode: false,
		MiniBye:     true,

		// Icons
		IconsVariant: "plain",

		// Reader
		ReaderName:          "",
		ReaderReadInBrowser: false,

		// History
		HistorySaveOnRead:     true,
		HistorySaveOnDownload: false,

		// Mangadex
		MangadexLanguage:                "en",
		MangadexNSFW:                    false,
		MangadexShowUnavailableChapters: false,

		// Logs
		LogsWrite: false,
		LogsLevel: "info",

		// Anilist
		AnilistEnable: false,
	}

	for field, value := range fields {
		viper.SetDefault(field, value)
	}
}

// resolveAliases resolves the aliases for the paths
func resolveAliases() {
	home := lo.Must(os.UserHomeDir())
	path := viper.GetString(DownloaderPath)

	switch runtime.GOOS {
	case "windows":
		path = strings.ReplaceAll(path, "%USERPROFILE%", home)
	case "darwin", "linux":
		path = strings.ReplaceAll(path, "$HOME", home)
		path = strings.ReplaceAll(path, "~", home)
	default:
		panic("unsupported OS: " + runtime.GOOS)
	}
}
