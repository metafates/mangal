package inline

import (
	"fmt"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/downloader"
	"github.com/metafates/mangal/source"
	"github.com/spf13/viper"
)

func Run(options *Options) error {
	mangas, err := options.Source.Search(options.Query)
	if err != nil {
		return err
	}

	// manga picker can only be none if json is set
	if options.MangaPicker.IsNone() {
		// preload all chapters
		for _, manga := range mangas {
			if err = jsonUpdateChapters(manga, options); err != nil {
				return err
			}
		}

		return printAsJson(mangas)
	}

	var chapters []*source.Chapter

	if len(mangas) == 0 {
		chapters = []*source.Chapter{}
	} else {
		manga := options.MangaPicker.Unwrap()(mangas)

		chapters, err = options.Source.ChaptersOf(manga)
		if err != nil {
			return err
		}

		chapters, err = options.ChaptersFilter(chapters)
		if err != nil {
			return err
		}

		if options.Json {
			if err = jsonUpdateChapters(manga, options); err != nil {
				return err
			}

			return printAsJson([]*source.Manga{manga})
		}
	}

	for _, chapter := range chapters {
		if options.Download {
			path, err := downloader.Download(chapter, func(string) {})
			if err != nil && viper.GetBool(constant.DownloaderStopOnError) {
				return err
			}

			fmt.Println(path)
		} else {
			err := downloader.Read(chapter, func(string) {})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
