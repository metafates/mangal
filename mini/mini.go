package mini

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/converter"
	"github.com/metafates/mangal/provider"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/util"
	"github.com/samber/lo"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
	"sync"
)

func Run(download bool) error {
	conv, err := converter.Get(viper.GetString(config.FormatsUse))
	if err != nil {
		return err
	}

	s, err := selectSource()
	if err != nil {
		return err
	}

	mangas, err := searchMangas(s)
	if err != nil {
		return err
	}

	manga, err := selectManga(mangas)
	if err != nil {
		return err
	}

	if !download {
		chapter, err := selectChapter(s, manga)
		if err != nil {
			return err
		}

		pages, err := s.PagesOf(chapter)
		if err != nil {
			return err
		}

		fmt.Printf("Downloading %d pages\n", len(pages))
		err = chapter.DownloadPages()
		if err != nil {
			return err
		}

		path, err := conv.SaveTemp(chapter)
		if err != nil {
			return err
		}

		err = open.Start(path)
		if err != nil {
			return err
		}
	} else {
		chapters, err := selectChapters(s, manga)
		if err != nil {
			return err
		}

		for _, chapter := range chapters {
			fmt.Printf("Downloading %s\n", chapter.Name)
			_, err = s.PagesOf(chapter)
			if err != nil {
				return err
			}

			err = chapter.DownloadPages()
			if err != nil {
				return err
			}
		}

		wg := sync.WaitGroup{}
		wg.Add(len(chapters))
		for _, chapter := range chapters {
			go func(chapter *source.Chapter) {
				defer wg.Done()

				if err != nil {
					return
				}

				_, err = conv.Save(chapter)
			}(chapter)
		}
		wg.Wait()

		if err != nil {
			return err
		}

		fmt.Println("Saved")
	}

	return nil
}

func selectSource() (source.Source, error) {

	defaultProviders := provider.DefaultProviders()
	customSources, err := source.AvailableCustomSources()

	if err != nil {
		return nil, err
	}

	var sources = make(map[string]func() (source.Source, error))

	for name, s := range customSources {
		sources[name] = func() (source.Source, error) {
			return source.LoadSource(s, true)
		}
	}

	for name, p := range defaultProviders {
		sources[name] = func() (source.Source, error) {
			return p.CreateSource(), nil
		}
	}

	options := lo.Keys(sources)
	slices.Sort(options)

	prompt := survey.Select{
		Message: "Select a source",
		Options: options,
		VimMode: viper.GetBool(config.MiniVimMode),
	}

	var sourceName string
	err = survey.AskOne(&prompt, &sourceName, survey.WithValidator(survey.Required))
	if err != nil {
		return nil, err
	}

	return sources[sourceName]()
}

func searchMangas(s source.Source) ([]*source.Manga, error) {
	prompt := survey.Input{
		Message: "Search manga",
	}

	var query string
	err := survey.AskOne(&prompt, &query, survey.WithValidator(survey.Required))
	if err != nil {
		return nil, err
	}

	mangas, err := s.Search(query)
	if err != nil {
		return nil, err
	}

	return mangas, nil
}

func selectManga(mangas []*source.Manga) (*source.Manga, error) {
	var m = make(map[string]*source.Manga)

	for _, manga := range mangas {
		m[util.PrettyTrim(manga.Name, 30)] = manga
	}

	options := lo.Keys(m)
	slices.SortFunc(options, func(a, b string) bool {
		return m[a].Index < m[b].Index
	})

	prompt := survey.Select{
		Message: "Select manga",
		Options: options,
		VimMode: viper.GetBool(config.MiniVimMode),
	}

	var mangaName string
	err := survey.AskOne(&prompt, &mangaName, survey.WithValidator(survey.Required))
	if err != nil {
		return nil, err
	}

	return m[mangaName], nil
}

func selectChapter(s source.Source, manga *source.Manga) (*source.Chapter, error) {
	chapters, err := s.ChaptersOf(manga)
	if err != nil {
		return nil, err
	}

	var c = make(map[string]*source.Chapter)
	for _, chapter := range chapters {
		c[util.PrettyTrim(chapter.Name, 30)] = chapter
	}

	options := lo.Keys(c)
	slices.SortFunc(options, func(a, b string) bool {
		return c[a].Index < c[b].Index
	})

	prompt := survey.Select{
		Message: "Select chapter",
		Options: options,
		VimMode: viper.GetBool(config.MiniVimMode),
	}

	var chapterName string
	err = survey.AskOne(&prompt, &chapterName, survey.WithValidator(survey.Required))
	if err != nil {
		return nil, err
	}

	return c[chapterName], nil
}

func selectChapters(s source.Source, manga *source.Manga) ([]*source.Chapter, error) {
	chapters, err := s.ChaptersOf(manga)
	if err != nil {
		return nil, err
	}

	var c = make(map[string]*source.Chapter)
	for _, chapter := range chapters {
		c[util.PrettyTrim(chapter.Name, 30)] = chapter
	}

	options := lo.Keys(c)
	slices.SortFunc(options, func(a, b string) bool {
		return c[a].Index < c[b].Index
	})

	prompt := survey.MultiSelect{
		Message: "Select chapters",
		Options: options,
		VimMode: viper.GetBool(config.MiniVimMode),
	}

	var chapterNames []string
	err = survey.AskOne(&prompt, &chapterNames, survey.WithValidator(survey.Required))
	if err != nil {
		return nil, err
	}

	var chaptersToDownload []*source.Chapter
	for _, chapterName := range chapterNames {
		chaptersToDownload = append(chaptersToDownload, c[chapterName])
	}

	slices.SortFunc(chaptersToDownload, func(a, b *source.Chapter) bool {
		return a.Index < b.Index
	})

	return chaptersToDownload, nil
}
