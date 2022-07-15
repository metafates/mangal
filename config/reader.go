package config

type ReaderConfig struct {
	UseCustomReader bool   `toml:"use_custom_reader"`
	CustomReader    string `toml:"custom_reader"`
}
