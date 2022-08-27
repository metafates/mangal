package config

import "github.com/metafates/mangal/constant"

var EnvExposed = []string{
	// downloader
	constant.DownloaderPath,
	constant.DownloaderChapterNameTemplate,
	constant.DownloaderCreateMangaDir,
	constant.DownloaderDefaultSource,

	// formats
	constant.FormatsUse,

	// reader
	constant.ReaderCBZ,
	constant.ReaderPDF,
	constant.ReaderZIP,
	constant.RaderPlain,

	// history
	constant.HistorySaveOnRead,
	constant.HistorySaveOnDownload,

	// Logs
	constant.LogsWrite,
	constant.LogsLevel,

	// Anilist
	constant.AnilistEnable,
	constant.AnilistID,
	constant.AnilistSecret,
	constant.AnilistCode,
}
