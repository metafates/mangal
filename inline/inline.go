package inline

import (
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/downloader"
	"github.com/metafates/mangal/log"
	"github.com/metafates/mangal/source"
	"github.com/spf13/viper"
	"os"
)

func Run(options *Options) error {
	if options.Out == nil {
		options.Out = os.Stdout
	}

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

		marshalled, err := asJson(mangas)
		if err != nil {
			return err
		}

		_, err = options.Out.Write(marshalled)
		return err
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

			marshalled, err := asJson([]*source.Manga{manga})
			if err != nil {
				return err
			}

			_, err = options.Out.Write(marshalled)
			return err
		}
	}

	for _, chapter := range chapters {
		if options.Download {
			path, err := downloader.Download(chapter, func(string) {})
			if err != nil && viper.GetBool(constant.DownloaderStopOnError) {
				return err
			}

			_, err = options.Out.Write([]byte(path + "\n"))
			if err != nil {
				log.Warn(err)
			}
		} else {
			err := downloader.Read(chapter, func(string) {})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
