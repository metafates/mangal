package mini

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/metafates/mangal/converter"
	"github.com/metafates/mangal/provider"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/util"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
	"strings"
)

func Run() error {
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

	path, err := converter.Converters()["cbz"].Save(chapter)
	if err != nil {
		return err
	}

	fmt.Println("Saved to ", path)

	return nil
}

func selectSource() (source.Source, error) {

	defaultProviders := provider.DefaultProviders()
	customSources, err := source.AvailableCustomSources()

	if err != nil {
		return nil, err
	}

	defaultProvidersNames := lo.Keys(defaultProviders)
	slices.Sort(defaultProvidersNames)

	customSourcesNames := lo.Keys(customSources)
	slices.Sort(customSourcesNames)

	items := append(
		lo.Map(defaultProvidersNames, func(name string, _ int) string {
			return "Builtin: " + name
		}),
		customSourcesNames...,
	)

	prompt := promptui.Select{
		Label:     "Select a source",
		Items:     items,
		IsVimMode: true,
		Searcher: func(input string, index int) bool {
			return strings.Contains(items[index], input)
		},
		StartInSearchMode: true,
		Size:              10,
	}

	i, name, err := prompt.Run()

	if err != nil {
		return nil, err
	}

	var s source.Source

	if i >= len(defaultProviders) {
		path := customSources[name]
		s, err = source.LoadSource(path, true)
		if err != nil {
			return nil, err
		}
	} else {
		s = defaultProviders[name].CreateSource()
	}

	return s, nil
}

func searchMangas(s source.Source) ([]*source.Manga, error) {
	prompt := promptui.Prompt{
		Label:       "Search",
		HideEntered: true,
	}

	query, err := prompt.Run()
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
	prompt := promptui.Select{
		Label: "Select a manga",
		Items: lo.Map(mangas, func(manga *source.Manga, _ int) string {
			return manga.Name
		}),
		Size: 20,
	}

	i, _, err := prompt.Run()

	if err != nil {
		return nil, err
	}

	return mangas[i], nil
}

func selectChapter(s source.Source, manga *source.Manga) (*source.Chapter, error) {
	chapters, err := s.ChaptersOf(manga)
	if err != nil {
		return nil, err
	}

	prompt := promptui.Select{
		Label: "Select a chapter",
		Items: lo.Map(chapters, func(chapter *source.Chapter, _ int) string {
			return util.PrettyTrim(chapter.Name, 30)
		}),
		Size: 20,
	}

	i, _, err := prompt.Run()

	if err != nil {
		return nil, err
	}

	return chapters[i], nil
}
