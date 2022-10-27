package inline

import (
	"encoding/json"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/source"
	"github.com/spf13/viper"
)

func asJson(manga []*source.Manga) (marshalled []byte, err error) {
	return json.Marshal(&struct {
		Manga []*source.Manga
	}{
		Manga: manga,
	})
}

func jsonUpdateChapters(manga *source.Manga, options *Options) error {
	var err error
	chapters, _ := options.Source.ChaptersOf(manga)
	if options.ChaptersFilter.IsPresent() {
		chapters, err = options.ChaptersFilter.MustGet()(chapters)
		if err != nil {
			return err
		}
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

	if viper.GetBool(constant.MetadataFetchAnilist) {
		_ = manga.PopulateMetadata(func(string) {})
	}

	return nil
}
