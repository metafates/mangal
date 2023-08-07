package config

import (
	"text/template"
	"time"

	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/icon"
	"github.com/mangalorg/mangal/nametemplate/util"
)

type config struct {
	Icons     field[string]
	Read      configRead
	Download  configDownload
	TUI       configTUI
	Providers configProviders
}

type configRead struct {
	Format         field[string]
	History        configReadHistory
	DownloadOnRead field[bool]
}

type configReadHistory struct {
	Anilist field[bool]
	Local   field[bool]
}

type configDownload struct {
	Format       field[string]
	Path         field[string]
	Strict       field[bool]
	SkipIfExists field[bool]
	Manga        configDownloadManga
	Volume       configDownloadVolume
	Chapter      configDownloadChapter
	Metadata     configDownloadMetadata
}

type configDownloadManga struct {
	CreateDir     field[bool]
	Cover, Banner field[bool]
	NameTemplate  field[string]
}

type configDownloadVolume struct {
	CreateDir    field[bool]
	NameTemplate field[string]
}

type configDownloadChapter struct {
	NameTemplate field[string]
}

type configDownloadMetadata struct {
	ComicInfoXML field[bool]
	SeriesJSON   field[bool]
}

type configTUI struct {
	ExpandSingleVolume field[bool]
}

type configProviders struct {
	Cache configProvidersCache
}

type configProvidersCache struct {
	TTL field[string]
}

var Config = config{
	Icons: register(field[string]{
		key:          "icons",
		defaultValue: icon.TypeASCII.String(),
		description:  "Icon format to use",
		init: func(s string) error {
			t, err := icon.TypeString(s)
			if err != nil {
				return err
			}

			icon.SetType(t)
			return nil
		},
	}),
	Read: configRead{
		Format: register(field[string]{
			key:          "read.format",
			defaultValue: libmangal.FormatPDF.String(),
			description:  "Format to read chapters in",
			init: func(s string) error {
				_, err := libmangal.FormatString(s)
				return err
			},
		}),
		History: configReadHistory{
			Anilist: register(field[bool]{
				key:          "read.history.anilist",
				defaultValue: true,
				description:  "Sync to Anilist reading history if logged in.",
			}),
			Local: register(field[bool]{
				key:          "read.history.local",
				defaultValue: true,
				description:  "Save to local history",
			}),
		},
		DownloadOnRead: register(field[bool]{
			key:          "read.download_on_read",
			defaultValue: false,
			description:  "Download chapter to the default directory when opening for reading",
		}),
	},
	Download: configDownload{
		Path: register(field[string]{
			key:          "download.path",
			defaultValue: ".",
			description:  "Path where chapters will be downloaded",
			transform:    expandPath,
		}),
		Format: register(field[string]{
			key:          "download.format",
			defaultValue: libmangal.FormatPDF.String(),
			description:  "Format to download chapters in",
			init: func(s string) error {
				_, err := libmangal.FormatString(s)
				return err
			},
		}),
		Strict: register(field[bool]{
			key:          "download.strict",
			defaultValue: false,
			description:  "If during metadata/banner/cover creation error occurs downloader will return it immediately and chapter won't be downloaded",
		}),
		SkipIfExists: register(field[bool]{
			key:          "download.skip_if_exists",
			defaultValue: true,
			description:  "Skip downloading chapter if its already downloaded (exists at path). Metadata will still be created if needed.",
		}),
		Manga: configDownloadManga{
			CreateDir: register(field[bool]{
				key:          "download.manga.create_dir",
				defaultValue: true,
				description:  "Create manga directory",
			}),
			Cover: register(field[bool]{
				key:          "download.manga.cover",
				defaultValue: false,
				description:  "Download manga cover",
			}),
			Banner: register(field[bool]{
				key:          "download.manga.banner",
				defaultValue: false,
				description:  "Download manga banner",
			}),
			NameTemplate: register(field[string]{
				key:          "download.manga.name_template",
				defaultValue: `{{ .Title | sanitize }}`,
				description:  "Template to use for naming downloaded mangas",
				init: func(s string) error {
					_, err := template.
						New("").
						Funcs(util.FuncMap).
						Parse(s)

					return err
				},
			}),
		},
		Volume: configDownloadVolume{
			CreateDir: register(field[bool]{
				key:          "download.volume.create_dir",
				defaultValue: false,
				description:  "Create volume directory",
			}),
			NameTemplate: register(field[string]{
				key:          "download.volume.name_template",
				defaultValue: `{{ printf "Vol. %d" .Number | sanitize }}`,
				description:  "Template to use for naming downloaded volumes",
				init: func(s string) error {
					_, err := template.
						New("").
						Funcs(util.FuncMap).
						Parse(s)

					return err
				},
			}),
		},
		Chapter: configDownloadChapter{
			NameTemplate: register(field[string]{
				key:          "download.chapter.name_template",
				defaultValue: `{{ printf "[%06.1f] %s" .Number .Title | sanitize }}`,
				description:  "Template to use for naming downloaded chapters",
				init: func(s string) error {
					_, err := template.
						New("").
						Funcs(util.FuncMap).
						Parse(s)

					return err
				},
			}),
		},
		Metadata: configDownloadMetadata{
			ComicInfoXML: register(field[bool]{
				key:          "download.metadata.comicinfo_xml",
				defaultValue: false,
				description:  "Generate `ComicInfo.xml` file",
			}),
			SeriesJSON: register(field[bool]{
				key:          "download.metadata.series_json",
				defaultValue: false,
				description:  "Generate `series.json` file",
			}),
		},
	},
	TUI: configTUI{
		ExpandSingleVolume: register(field[bool]{
			key:          "tui.expand_single_volume",
			defaultValue: true,
			description:  "Skip selecting volume if there's only one",
		}),
	},
	Providers: configProviders{
		Cache: configProvidersCache{
			TTL: register(field[string]{
				key:          "providers.cache.ttl",
				defaultValue: "24h",
				description:  `Time to live. A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".`,
				init: func(s string) error {
					_, err := time.ParseDuration(s)
					return err
				},
			}),
		},
	},
}
