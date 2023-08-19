package config

import (
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/icon"
)

type config struct {
	Icons     *registered[string, icon.Type]
	CLI       configCLI
	Read      configRead
	Download  configDownload
	TUI       configTUI
	Providers configProviders
	Library   configLibrary
}

type configCLI struct {
	ColoredHelp *registered[bool, bool]
	Mode        configCLIMode
}

type configCLIMode struct {
	Default *registered[string, Mode]
}

type configRead struct {
	Format         *registered[string, libmangal.Format]
	History        configReadHistory
	DownloadOnRead *registered[bool, bool]
}

type configReadHistory struct {
	Anilist *registered[bool, bool]
	Local   *registered[bool, bool]
}

type configDownload struct {
	Format       *registered[string, libmangal.Format]
	Path         *registered[string, string]
	Strict       *registered[bool, bool]
	SkipIfExists *registered[bool, bool]
	Manga        configDownloadManga
	Volume       configDownloadVolume
	Chapter      configDownloadChapter
	Metadata     configDownloadMetadata
}

type configDownloadManga struct {
	CreateDir     *registered[bool, bool]
	Cover, Banner *registered[bool, bool]
	NameTemplate  *registered[string, string]
}

type configDownloadVolume struct {
	CreateDir    *registered[bool, bool]
	NameTemplate *registered[string, string]
}

type configDownloadChapter struct {
	NameTemplate *registered[string, string]
}

type configDownloadMetadata struct {
	ComicInfoXML *registered[bool, bool]
	SeriesJSON   *registered[bool, bool]
}

type configTUI struct {
	ExpandSingleVolume *registered[bool, bool]
}

type configProviders struct {
	Cache configProvidersCache
}

type configProvidersCache struct {
	TTL *registered[string, string]
}

type configLibrary struct {
	Path *registered[string, string]
}
