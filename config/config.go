package config

import (
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/filesystem"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
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
	viper.AddConfigPath(Path())
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

func init() {
	paths := []string{
		Path(),
		SourcesPath(),
		LogsPath(),
	}

	for _, path := range paths {
		lo.Must0(filesystem.Get().MkdirAll(path, os.ModePerm))
	}
}

func Path() string {
	var path string

	customDir, present := os.LookupEnv("MANGAL_CONFIG_DIR")
	if present {
		path = customDir
	} else {
		path = filepath.Join(lo.Must(os.UserConfigDir()), constant.Mangal)
	}

	return path
}

func SourcesPath() string {
	return filepath.Join(Path(), "sources")
}

func LogsPath() string {
	return filepath.Join(Path(), "logs")
}
