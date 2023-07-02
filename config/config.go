package config

import (
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/icon"
)

type config struct {
	Icons    field[string]
	Read     configRead
	Download configDownload
}

type configRead struct {
	Format    field[string]
	Incognito field[bool]
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
	Dir           field[bool]
	Cover, Banner field[bool]
}

type configDownloadVolume struct {
	Dir field[bool]
}

type configDownloadMetadata struct {
	ComicInfoXML field[bool]
	SeriesJSON   field[bool]
}

var Config = config{
	Icons: register(field[string]{
		key:          "icon.type",
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
			Dir: register(field[bool]{
				key:          "download.manga.dir",
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
			Dir: register(field[bool]{
				key:          "download.volume.dir",
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
}
