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
	envPath, defined := os.LookupEnv(constant.Mangal + "_CONFIG_PATH")
	if defined {
		viper.AddConfigPath(envPath)
	}

	paths := lo.Must(Paths())

	for _, path := range paths {
		viper.AddConfigPath(path)
	}
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
	configDir := lo.Must(os.UserConfigDir())

	fields := map[string]any{
		// Downloader
		DownloaderPath:                ".",
		DownloaderChapterNameTemplate: "[{padded-index}] {chapter}",

		// Formats
		FormatsUse: "plain",

		// Sources
		SourcesPath: filepath.Join(configDir, constant.Mangal, "sources"),

		// Mini-mode
		MiniVimMode: false,
		MiniBye:     true,

		// Icons
		IconsVariant: "emoji",

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
	}

	for field, value := range fields {
		viper.SetDefault(field, value)
	}
}

// resolveAliases resolves the aliases for the downloader path
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

	viper.Set(DownloaderPath, path)
}

// Paths returns the paths to the config files
func Paths() ([]string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return []string{
		homeDir,
		filepath.Join(configDir, constant.Mangal),
	}, nil
}
