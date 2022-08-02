package config

const (
	DownloaderPath                = "downloader.path"
	DownloaderChapterNameTemplate = "downloader.chapter_name_template"
)

const (
	FormatsDefault = "formats.default"
)

const (
	ReaderName = "reader.name"
)

const (
	SourcesPath = "sources.path"
)

var envFields = []string{
	DownloaderPath,
	DownloaderChapterNameTemplate,
	FormatsDefault,
	ReaderName,
	SourcesPath,
}
