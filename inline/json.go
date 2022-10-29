package inline

import (
	"encoding/json"
	"github.com/metafates/mangal/anilist"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/source"
	"github.com/spf13/viper"
)

func asJson(manga []*source.Manga, options *Options) (marshalled []byte, err error) {
	type inlineManga struct {
		*source.Manga
		Anilist *anilist.Manga `json:"anilist"`
	}

	var m = make([]*inlineManga, len(manga))
	for i, manga := range manga {
		al := manga.Anilist.OrElse(nil)
		if !options.IncludeAnilistManga {
			al = nil
		}

		m[i] = &inlineManga{
			Manga:   manga,
			Anilist: al,
		}
	}

	return json.Marshal(&struct {
		Source string         `json:"source"`
		Query  string         `json:"query"`
		Result []*inlineManga `json:"result"`
	}{
		Result: m,
		Source: options.Source.Name(),
		Query:  options.Query,
	})
}

func prepareManga(manga *source.Manga, options *Options) error {
	var err error

	if options.IncludeAnilistManga {
		err = manga.BindWithAnilist()
		if err != nil {
			return err
		}
	}

	if options.ChaptersFilter.IsPresent() {
		chapters, err := options.Source.ChaptersOf(manga)
		if err != nil {
			return err
		}

		chapters, err = options.ChaptersFilter.MustGet()(chapters)
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
	} else {
		// clear chapters in case they were loaded from cache or something
		manga.Chapters = make([]*source.Chapter, 0)
	}

	if viper.GetBool(constant.MetadataFetchAnilist) {
		_ = manga.PopulateMetadata(func(string) {})
	}

	return nil
}
