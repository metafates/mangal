package config

import "github.com/metafates/mangal/constant"

var EnvExposed = []string{
	// downloader
	constant.DownloaderPath,
	constant.DownloaderChapterNameTemplate,
	constant.DownloaderCreateMangaDir,
	constant.DownloaderDefaultSource,
	constant.DownloaderAsync,
	constant.DownloaderDownloadCover,
	constant.DownloaderCreateMangaDir,
	constant.DownloaderCreateVolumeDir,
	constant.DownloaderStopOnError,

	// formats
	constant.FormatsUse,
	constant.FormatsSkipUnsupportedImages,

	// reader
	constant.ReaderCBZ,
	constant.ReaderPDF,
	constant.ReaderZIP,
	constant.RaderPlain,
	constant.ReaderReadInBrowser,

	// history
	constant.HistorySaveOnRead,
	constant.HistorySaveOnDownload,

	// metadata
	constant.MetadataFetchAnilist,
	constant.MetadataComicInfoXML,
	constant.MetadataSeriesJSON,

	// Logs
	constant.LogsWrite,
	constant.LogsLevel,
	constant.LogsJson,

	// Anilist
	constant.AnilistEnable,
	constant.AnilistID,
	constant.AnilistSecret,
	constant.AnilistCode,
}
