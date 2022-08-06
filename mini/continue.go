package mini

import (
	"errors"
	"github.com/AlecAivazis/survey/v2"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/history"
	"github.com/metafates/mangal/provider"
	"github.com/metafates/mangal/source"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

func continueReading() error {

	saved, err := history.Get()
	if err != nil {
		return err
	}

	if len(saved) == 0 {
		return errors.New("history is empty")
	}

	options := lo.Keys(saved)
	slices.Sort(options)

	prompt := survey.Select{
		Message:  "Select a manga",
		Options:  options,
		VimMode:  viper.GetBool(config.MiniVimMode),
		PageSize: pageSize,
	}

	var mangaName string
	err = survey.AskOne(&prompt, &mangaName, survey.WithValidator(survey.Required))
	if err != nil {
		return err
	}

	defaultProviders := provider.DefaultProviders()
	customProviders, err := provider.CustomProviders()

	if err != nil {
		return err
	}

	var sources = make(map[string]func() (source.Source, error))

	for _, p := range customProviders {
		sources[p.ID] = p.CreateSource
	}

	for _, p := range defaultProviders {
		sources[p.ID] = p.CreateSource
	}

	chap := saved[mangaName]
	s, err := sources[chap.SourceID]()
	if err != nil {
		return err
	}

	manga := &source.Manga{
		Name:     chap.MangaName,
		URL:      chap.MangaURL,
		ID:       chap.MangaID,
		Index:    0,
		SourceID: chap.SourceID,
		Chapters: make([]*source.Chapter, 0),
	}

	chapters, err := s.ChaptersOf(manga)
	if err != nil {
		return err
	}

	_, index, ok := lo.FindIndexOf(chapters, func(c *source.Chapter) bool {
		return c.URL == chap.URL
	})

	if !ok {
		return errors.New("chapter not found")
	}

	chapter, err := selectChapter(chapters, index)

	if err != nil {
		return err
	}

	return readChapter(s, chapter)
}
