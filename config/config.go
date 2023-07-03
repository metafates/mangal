package config

import (
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/icon"
)

type config struct {
	Icons    field[string]
	Read     configRead
	Download configDownload
	TUI      configTUI
}

type configRead struct {
	Format         field[string]
	Incognito      field[bool]
	DownloadOnRead field[bool]
}

type configDownload struct {
	Format       field[string]
	Path         field[string]
	Strict       field[bool]
	SkipIfExists field[bool]
	Manga        configDownloadManga
	Volume       configDownloadVolume
	Metadata     configDownloadMetadata
}

type configDownloadManga struct {
	CreateDir     field[bool]
	Cover, Banner field[bool]
}

type configDownloadVolume struct {
	CreateDir field[bool]
}

type configDownloadMetadata struct {
	ComicInfoXML field[bool]
	SeriesJSON   field[bool]
}

type configTUI struct {
	ExpandSingleVolume field[bool]
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
		}),
		Incognito: register(field[bool]{
			key:          "read.incognito",
			defaultValue: false,
			description:  "Won't sync to Anilist reading history if logged in.",
		}),
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
		},
		Volume: configDownloadVolume{
			CreateDir: register(field[bool]{
				key:          "download.volume.create_dir",
				defaultValue: false,
				description:  "Create volume directory",
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
}
