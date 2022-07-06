package config

type DownloaderConfig struct {
	ChapterNameTemplate string `toml:"chapter_name_template"`
	Path                string `toml:"path"`
}
