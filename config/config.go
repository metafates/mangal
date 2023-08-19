package config

import (
	"text/template"
	"time"

	"github.com/adrg/xdg"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/icon"
	"github.com/mangalorg/mangal/nametemplate/util"
)

var Config = config{
	Icons: reg(Field[string, icon.Type]{
		Key:         "icons",
		Default:     icon.TypeASCII,
		Description: "Icon format to use",
		Unmarshal: func(s string) (icon.Type, error) {
			return icon.TypeString(s)
		},
		Marshal: func(i icon.Type) (string, error) {
			return i.String(), nil
		},
	}),
	CLI: configCLI{
		ColoredHelp: reg(Field[bool, bool]{
			Key:         "cli.colored_help",
			Default:     true,
			Description: "Enable colors in cli help",
		}),
		Mode: configCLIMode{
			Default: reg(Field[string, Mode]{
				Key:         "cli.mode.default",
				Default:     ModeTUI,
				Description: "Default mode to use when no subcommand is given",
				Unmarshal: func(s string) (Mode, error) {
					return ModeString(s)
				},
				Marshal: func(mode Mode) (string, error) {
					return mode.String(), nil
				},
			}),
		},
	},
	Read: configRead{
		Format: reg(Field[string, libmangal.Format]{
			Key:         "read.format",
			Default:     libmangal.FormatPDF,
			Description: "Format to read chapters in",
			Unmarshal: func(s string) (libmangal.Format, error) {
				return libmangal.FormatString(s)
			},
			Marshal: func(format libmangal.Format) (string, error) {
				return format.String(), nil
			},
		}),
		History: configReadHistory{
			Anilist: reg(Field[bool, bool]{
				Key:         "read.history.anilist",
				Default:     true,
				Description: "Sync to Anilist reading history if logged in.",
			}),
			Local: reg(Field[bool, bool]{
				Key:         "read.history.local",
				Default:     true,
				Description: "Save to local history",
			}),
		},
		DownloadOnRead: reg(Field[bool, bool]{
			Key:         "read.download_on_read",
			Default:     false,
			Description: "Download chapter to the default directory when opening for reading",
		}),
	},
	Download: configDownload{
		Path: reg(Field[string, string]{
			Key:         "download.path",
			Default:     xdg.UserDirs.Download,
			Description: "Path where chapters will be downloaded",
			Unmarshal: func(s string) (string, error) {
				return expandPath(s)
			},
		}),
		Format: reg(Field[string, libmangal.Format]{
			Key:         "download.format",
			Default:     libmangal.FormatPDF,
			Description: "Format to download chapters in",
			Unmarshal: func(s string) (libmangal.Format, error) {
				return libmangal.FormatString(s)
			},
			Marshal: func(format libmangal.Format) (string, error) {
				return format.String(), nil
			},
		}),
		Strict: reg(Field[bool, bool]{
			Key:         "download.strict",
			Default:     false,
			Description: "If during metadata/banner/cover creation error occurs downloader will return it immediately and chapter won't be downloaded",
		}),
		SkipIfExists: reg(Field[bool, bool]{
			Key:         "download.skip_if_exists",
			Default:     true,
			Description: "Skip downloading chapter if its already downloaded (exists at path). Metadata will still be created if needed.",
		}),
		Manga: configDownloadManga{
			CreateDir: reg(Field[bool, bool]{
				Key:         "download.manga.create_dir",
				Default:     true,
				Description: "Create manga directory",
			}),
			Cover: reg(Field[bool, bool]{
				Key:         "download.manga.cover",
				Default:     false,
				Description: "Download manga cover",
			}),
			Banner: reg(Field[bool, bool]{
				Key:         "download.manga.banner",
				Default:     false,
				Description: "Download manga banner",
			}),
			NameTemplate: reg(Field[string, string]{
				Key:         "download.manga.name_template",
				Default:     `{{ .Title | sanitize }}`,
				Description: "Template to use for naming downloaded mangas",
				Validate: func(s string) error {
					_, err := template.
						New("").
						Funcs(util.FuncMap).
						Parse(s)

					return err
				},
			}),
		},
		Volume: configDownloadVolume{
			CreateDir: reg(Field[bool, bool]{
				Key:         "download.volume.create_dir",
				Default:     false,
				Description: "Create volume directory",
			}),
			NameTemplate: reg(Field[string, string]{
				Key:         "download.volume.name_template",
				Default:     `{{ printf "Vol. %d" .Number | sanitize }}`,
				Description: "Template to use for naming downloaded volumes",
				Validate: func(s string) error {
					_, err := template.
						New("").
						Funcs(util.FuncMap).
						Parse(s)

					return err
				},
			}),
		},
		Chapter: configDownloadChapter{
			NameTemplate: reg(Field[string, string]{
				Key:         "download.chapter.name_template",
				Default:     `{{ printf "[%06.1f] %s" .Number .Title | sanitize }}`,
				Description: "Template to use for naming downloaded chapters",
				Validate: func(s string) error {
					_, err := template.
						New("").
						Funcs(util.FuncMap).
						Parse(s)

					return err
				},
			}),
		},
		Metadata: configDownloadMetadata{
			ComicInfoXML: reg(Field[bool, bool]{
				Key:         "download.metadata.comicinfo_xml",
				Default:     false,
				Description: "Generate `ComicInfo.xml` file",
			}),
			SeriesJSON: reg(Field[bool, bool]{
				Key:         "download.metadata.series_json",
				Default:     false,
				Description: "Generate `series.json` file",
			}),
		},
	},
	TUI: configTUI{
		ExpandSingleVolume: reg(Field[bool, bool]{
			Key:         "tui.expand_single_volume",
			Default:     true,
			Description: "Skip selecting volume if there's only one",
		}),
	},
	Providers: configProviders{
		Cache: configProvidersCache{
			TTL: reg(Field[string, string]{
				Key:         "providers.cache.ttl",
				Default:     "24h",
				Description: `Time to live. A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".`,
				Validate: func(s string) error {
					_, err := time.ParseDuration(s)
					return err
				},
			}),
		},
	},
	Library: configLibrary{
		Path: reg(Field[string, string]{
			Key:         "library.path",
			Default:     "",
			Description: "Path to the manga library. Empty string will fallback to the download.path",
		}),
	},
}
