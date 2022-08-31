package inline

import (
	"encoding/json"
	"fmt"
	"github.com/metafates/mangal/source"
	"os"
)

func printAsJson(manga []*source.Manga) error {
	marshalled, err := json.Marshal(&struct {
		Manga []*source.Manga
	}{
		Manga: manga,
	})

	if err != nil {
		return err
	}

	_, err = fmt.Fprint(os.Stdout, string(marshalled))
	return err
}

func jsonUpdateChapters(manga *source.Manga, options *Options) error {
	chapters, _ := options.Source.ChaptersOf(manga)
	chapters, err := options.ChaptersFilter(chapters)

	if err != nil {
		return err
	}

	manga.Chapters = chapters

	if options.PopulatePages {
		for _, chapter := range chapters {
			_, err := options.Source.PagesOf(chapter)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
