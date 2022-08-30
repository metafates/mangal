package inline

import (
	"encoding/json"
	"errors"
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

	if len(mangas) == 0 {
		return errors.New("no mangas found")
	}

	if options.Json && options.All {
		// pre-load _all_ chapters
		for _, manga := range mangas {
			if err = jsonUpdateChapters(manga, options); err != nil {
				return err
			}
		}

		return jsonPrint(&JsonData{Manga: mangas})
	}

	manga := options.MangaPicker(mangas)

	chapters, err := options.Source.ChaptersOf(manga)
	if err != nil {
		return err
	}

	if len(chapters) == 0 {
		return errors.New("no chapters found")
	}

	chapters, err = options.ChapterFilter(chapters)
	if err != nil {
		return err
	}

	if options.Json {
		if err = jsonUpdateChapters(manga, options); err != nil {
			return err
		}

		mangas := make([]*source.Manga, 1)
		mangas[0] = manga
		return jsonPrint(&JsonData{Manga: mangas})
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

func jsonUpdateChapters(manga *source.Manga, options *Options) error {
	chapters, _ := options.Source.ChaptersOf(manga)
	chapters, err := options.ChapterFilter(chapters)

	if err == nil {
		manga.Chapters = chapters
	}

	return err
}

func jsonPrint(data *JsonData) error {
	marshalledManga, err := json.Marshal(data)

	if err != nil {
		return err
	}

	fmt.Print(string(marshalledManga))

	return nil
}
