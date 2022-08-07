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

// resolveAliases resolves the aliases for the paths
func resolveAliases() {
	home := lo.Must(os.UserHomeDir())
	path := viper.GetString(DownloaderPath)
	srcPath := viper.GetString(SourcesPath)

	switch runtime.GOOS {
	case "windows":
		path = strings.ReplaceAll(path, "%USERPROFILE%", home)
		srcPath = strings.ReplaceAll(srcPath, "%USERPROFILE%", home)
	case "darwin", "linux":
		path = strings.ReplaceAll(path, "$HOME", home)
		srcPath = strings.ReplaceAll(srcPath, "$HOME", home)
		path = strings.ReplaceAll(path, "~", home)
		srcPath = strings.ReplaceAll(srcPath, "~", home)
	default:
		panic("unsupported OS: " + runtime.GOOS)
	}

	viper.Set(DownloaderPath, path)
	viper.Set(SourcesPath, srcPath)
}

// Paths returns the paths to the config files
func Paths() ([]string, error) {
	var paths []string

	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	paths = append(paths, filepath.Join(configDir, constant.Mangal))

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	paths = append(paths, homeDir)

	envPath, defined := os.LookupEnv(strings.ToUpper(constant.Mangal) + "_CONFIG_PATH")
	if defined {
		paths = append(paths, envPath)
	}

	return paths, nil
}
