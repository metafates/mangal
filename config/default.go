package config

import "github.com/metafates/mangal/constant"

var defaultValues = map[string]any{
	// Downloader
	constant.DownloaderPath:                ".",
	constant.DownloaderChapterNameTemplate: "[{padded-index}] {chapter}",
	constant.DownloaderAsync:               true,
	constant.DownloaderCreateMangaDir:      true,
	constant.DownloaderCreateVolumeDir:     false,
	constant.DownloaderDefaultSource:       "",
	constant.DownloaderStopOnError:         false,
	constant.DownloaderDownloadCover:       false,

	// Formats
	constant.FormatsUse:                   "pdf",
	constant.FormatsSkipUnsupportedImages: true,
	constant.FormatsFetchAnilistMetadata:  true,

	// Mini-mode
	constant.MiniSearchLimit: 20,

	// Icons
	constant.IconsVariant: "plain",

	// Reader
	constant.ReaderPDF:           "",
	constant.ReaderCBZ:           "",
	constant.ReaderZIP:           "",
	constant.RaderPlain:          "",
	constant.ReaderReadInBrowser: false,

	// History
	constant.HistorySaveOnRead:     true,
	constant.HistorySaveOnDownload: false,

	// Mangadex
	constant.MangadexLanguage:                "en",
	constant.MangadexNSFW:                    false,
	constant.MangadexShowUnavailableChapters: false,

	// Installer
	constant.InstallerUser:   "metafates",
	constant.InstallerRepo:   "mangal-scrapers",
	constant.InstallerBranch: "main",

	// Gen
	constant.GenAuthor: "",

	// Logs
	constant.LogsWrite: false,
	constant.LogsLevel: "info",
	constant.LogsJson:  false,

	// Anilist
	constant.AnilistEnable: false,
}
