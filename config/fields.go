package config

const (
	DownloaderPath                = "downloader.path"
	DownloaderChapterNameTemplate = "downloader.chapter_name_template"
)

const (
	FormatsUse = "formats.use"
)

const (
	ReaderName          = "reader.name"
	ReaderReadInBrowser = "reader.read_in_browser"
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

const (
	MangadexLanguage                = "mangadex.language"
	MangadexNSFW                    = "mangadex.nsfw"
	MangadexShowUnavailableChapters = "mangadex.show_unavailable_chapters"
)

const (
	AnilistEnable = "anilist.enable"
	AnilistID     = "anilist.id"
	AnilistSecret = "anilist.secret"
	AnilistCode   = "anilist.code"
)

const (
	LogsWrite = "logs.write"
	LogsPath  = "logs.path"
	LogsLevel = "logs.level"
)

var EnvExposed = []string{
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

	// Logs
	LogsWrite,
	LogsPath,
	LogsLevel,

	// Anilist
	AnilistEnable,
	AnilistID,
	AnilistSecret,
	AnilistCode,
}
