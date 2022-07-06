package common

// Environment variables
const (
	// EnvConfigPath points to config directory
	EnvConfigPath = "MANGAL_CONFIG_PATH"

	// EnvDownloadPath points to download directory
	EnvDownloadPath = "MANGAL_DOWNLOAD_PATH"

	// EnvDefaultFormat defines default format
	EnvDefaultFormat = "MANGAL_DEFAULT_FORMAT"

	// EnvCustomReader defines custom reader
	EnvCustomReader = "MANGAL_CUSTOM_READER"
)

var AvailableEnvVars = map[string]string{
	EnvConfigPath:    "Points to the config directory",
	EnvDownloadPath:  "Points to the downloads directory",
	EnvDefaultFormat: "Defines default format",
	EnvCustomReader:  "Defines custom reader",
}
