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
	HistorySaveOnRead     = "history.save_on_read"
	HistorySaveOnDownload = "history.save_on_download"
)

const (
	SourcesPath = "sources.path"
)

const (
	MiniVimMode = "mini.vim_mode"
	MiniBye     = "mini.bye"
)

const (
	IconsVariant = "icons.variant"
)

var envFields = []string{
	// downloader
	DownloaderPath,
	DownloaderChapterNameTemplate,

	// formats
	FormatsUse,

	// reader
	ReaderName,

	// history
	HistorySaveOnRead,
	HistorySaveOnDownload,

	// sources
	SourcesPath,

	// mini
	MiniVimMode,
	MiniBye,

	// Icons
	IconsVariant,
}
