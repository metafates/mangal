package inline

import (
	"errors"
	"fmt"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/downloader"
	"github.com/spf13/viper"
)

func Run(options *Options) error {
	mangas, err := options.Source.Search(options.Query)
	if err != nil {
		return err
	}

	if len(mangas) == 0 {
		return errors.New("no mangas found")
	}

	manga := options.MangaPicker(mangas)

	chapters, err := options.Source.ChaptersOf(manga)
	if err != nil {
		return err
	}

	if len(chapters) == 0 {
		return errors.New("no chapters found")
	}

	chapters = options.ChapterFilter(chapters)

	for _, chapter := range chapters {
		if options.Download {
			path, err := downloader.Download(options.Source, chapter, func(string) {})
			if err != nil && viper.GetBool(constant.DownloaderStopOnError) {
				return err
			}

			fmt.Println(path)
		} else {
			err := downloader.Read(options.Source, chapter, func(string) {})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
