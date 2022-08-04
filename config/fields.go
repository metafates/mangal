package config

const (
	DownloaderPath                = "downloader.path"
	DownloaderChapterNameTemplate = "downloader.chapter_name_template"
)

const (
	FormatsUse = "formats.use"
)

const (
	ReaderName = "reader.name"
)

const (
	SourcesPath = "sources.path"
)

const (
	MiniVimMode = "mini.vim_mode"
)

const (
	IconsVariant = "icons.variant"
)

var envFields = []string{
	DownloaderPath,
	DownloaderChapterNameTemplate,
	FormatsUse,
	ReaderName,
	SourcesPath,
}
