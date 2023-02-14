package key

// DefinedFieldsCount is the number of fields defined in this package.
// You have to manually update this number when you add a new field
// to check later if every field has a defined default value
const DefinedFieldsCount = 53

const (
	DownloaderPath                = "downloader.path"
	DownloaderChapterNameTemplate = "downloader.chapter_name_template"
	DownloaderAsync               = "downloader.async"
	DownloaderCreateMangaDir      = "downloader.create_manga_dir"
	DownloaderCreateVolumeDir     = "downloader.create_volume_dir"
	DownloaderDefaultSources      = "downloader.default_sources"
	DownloaderStopOnError         = "downloader.stop_on_error"
	DownloaderDownloadCover       = "downloader.download_cover"
	DownloaderRedownloadExisting  = "downloader.redownload_existing"
	DownloaderReadDownloaded      = "downloader.read_downloaded"
)

const (
	FormatsUse                   = "formats.use"
	FormatsSkipUnsupportedImages = "formats.skip_unsupported_images"
)

const (
	MetadataFetchAnilist                      = "metadata.fetch_anilist"
	MetadataComicInfoXML                      = "metadata.comic_info_xml"
	MetadataComicInfoXMLAddDate               = "metadata.comic_info_xml_add_date"
	MetadataComicInfoXMLAlternativeDate       = "metadata.comic_info_xml_alternative_date"
	MetadataComicInfoXMLTagRelevanceThreshold = "metadata.comic_info_xml_tag_relevance_threshold"
	MetadataSeriesJSON                        = "metadata.series_json"
)

const (
	ReaderPDF           = "reader.pdf"
	ReaderCBZ           = "reader.cbz"
	ReaderZIP           = "reader.zip"
	RaderPlain          = "reader.plain"
	ReaderBrowser       = "reader.browser"
	ReaderFolder        = "reader.folder"
	ReaderReadInBrowser = "reader.read_in_browser"
)

const (
	HistorySaveOnRead     = "history.save_on_read"
	HistorySaveOnDownload = "history.save_on_download"
)

const (
	SearchShowQuerySuggestions = "search.show_query_suggestions"
)

const (
	MiniSearchLimit = "mini.search_limit"
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
	AnilistEnable            = "anilist.enable"
	AnilistID                = "anilist.id"
	AnilistSecret            = "anilist.secret"
	AnilistCode              = "anilist.code"
	AnilistLinkOnMangaSelect = "anilist.link_on_manga_select"
)

const (
	TUIItemSpacing        = "tui.item_spacing"
	TUIReadOnEnter        = "tui.read_on_enter"
	TUISearchPromptString = "tui.search_prompt"
	TUIShowURLs           = "tui.show_urls"
	TUIShowDownloadedPath = "tui.show_downloaded_path"
	TUIReverseChapters    = "tui.reverse_chapters"
)

const (
	InstallerUser   = "installer.user"
	InstallerRepo   = "installer.repo"
	InstallerBranch = "installer.branch"
)

const (
	GenAuthor = "gen.author"
)

const (
	LogsWrite = "logs.write"
	LogsLevel = "logs.level"
	LogsJson  = "logs.json"
)

const (
	CliColored      = "cli.colored"
	CliVersionCheck = "cli.version_check"
)
