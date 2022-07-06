package config

type UIConfig struct {
	Fullscreen          bool
	Prompt              string
	Title               string
	Placeholder         string
	Mark                string
	ChapterNameTemplate string `toml:"chapter_name_template"`
}
