package mini

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/converter"
	"github.com/metafates/mangal/provider"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/style"
	"github.com/metafates/mangal/util"
	"github.com/samber/lo"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
	"os"
	"time"
)

var (
	pageSize = 15
	trimAt   = 30
)

func Run(download bool) error {
	if w, _, err := util.TerminalSize(); err == nil {
		trimAt = lo.Max([]int{trimAt, w - 10})
	}

	conv, err := converter.Get(viper.GetString(config.FormatsUse))
	if err != nil {
		return err
	}

	src, err := selectSource()
	if err != nil {
		return err
	}

	mangas, err := searchMangas(src)
	if err != nil {
		return err
	}

	manga, err := selectManga(mangas)
	if err != nil {
		return err
	}

	if !download {
		chapter, err := selectChapter(src, manga)
		if err != nil {
			return err
		}

		pages, err := src.PagesOf(chapter)
		if err != nil {
			return err
		}

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
		s.Suffix = fmt.Sprintf(" Downloading %d pages", len(pages))
		s.FinalMSG = style.Combined(style.Padding(1), style.Magenta)("ฅ^•ﻌ•^ฅ\nDone! Bye")
		lo.Must0(s.Color("bold", "magenta"))
		s.Start()
		err = chapter.DownloadPages()
		if err != nil {
			return err
		}

		s.Suffix = " Converting pages"
		path, err := conv.SaveTemp(chapter)
		if err != nil {
			return err
		}

		if reader := viper.GetString(config.ReaderName); reader != "" {
			s.Suffix = fmt.Sprintf(" Opening \"%s\"", reader)
			err = open.StartWith(path, reader)
		} else {
			s.Suffix = " Opening"
			err = open.Start(path)
		}

		if err != nil {
			return err
		}

		s.Stop()
	} else {
		chapters, err := selectChapters(src, manga)
		if err != nil {
			return err
		}

		var counter int
		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
		lo.Must0(s.Color("bold", "magenta"))
		s.Suffix = " Starting..."
		s.FinalMSG = style.Combined(style.Padding(1), style.Magenta)("ฅ^•ﻌ•^ฅ\nDone! Bye")
		s.Start()

		for _, chapter := range chapters {
			counter++

			s.Suffix = fmt.Sprintf(" [%d/%d] Getting pages of %s", counter, len(chapters), style.Trim(40)(chapter.Name))
			_, err = src.PagesOf(chapter)
			if err != nil {
				return err
			}

			s.Suffix = fmt.Sprintf(" [%d/%d] Downloading %d pages", counter, len(chapters), len(chapter.Pages))
			err = chapter.DownloadPages()
			if err != nil {
				return err
			}

			s.Suffix = fmt.Sprintf(" [%d/%d] Converting to %s", counter, len(chapters), viper.GetString(config.FormatsUse))
			_, err = conv.Save(chapter)
		}

		s.Stop()
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
		Message:  "Select a source",
		Options:  options,
		VimMode:  viper.GetBool(config.MiniVimMode),
		PageSize: pageSize,
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
		m[style.Trim(trimAt)(manga.Name)] = manga
	}

	options := lo.Keys(m)
	slices.SortFunc(options, func(a, b string) bool {
		return m[a].Index < m[b].Index
	})

	prompt := survey.Select{
		Message:  "Select manga",
		Options:  options,
		VimMode:  viper.GetBool(config.MiniVimMode),
		PageSize: pageSize,
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
		c[style.Trim(trimAt)(chapter.Name)] = chapter
	}

	options := lo.Keys(c)
	slices.SortFunc(options, func(a, b string) bool {
		return c[a].Index < c[b].Index
	})

	prompt := survey.Select{
		Message:  "Select chapter",
		Options:  options,
		VimMode:  viper.GetBool(config.MiniVimMode),
		PageSize: pageSize,
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
		c[style.Trim(trimAt)(chapter.Name)] = chapter
	}

	options := lo.Keys(c)
	slices.SortFunc(options, func(a, b string) bool {
		return c[a].Index < c[b].Index
	})

	// Remove selection answer
	survey.MultiSelectQuestionTemplate = `
{{- define "option"}}
    {{- if eq .SelectedIndex .CurrentIndex }}{{color .Config.Icons.SelectFocus.Format }}{{ .Config.Icons.SelectFocus.Text }}{{color "reset"}}{{else}} {{end}}
    {{- if index .Checked .CurrentOpt.Index }}{{color .Config.Icons.MarkedOption.Format }} {{ .Config.Icons.MarkedOption.Text }} {{else}}{{color .Config.Icons.UnmarkedOption.Format }} {{ .Config.Icons.UnmarkedOption.Text }} {{end}}
    {{- color "reset"}}
    {{- " "}}{{- .CurrentOpt.Value}}
{{end}}
{{- if .ShowHelp }}{{- color .Config.Icons.Help.Format }}{{ .Config.Icons.Help.Text }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color .Config.Icons.Question.Format }}{{ .Config.Icons.Question.Text }} {{color "reset"}}
{{- color "default+hb"}}{{ .Message }}{{ .FilterMessage }}{{color "reset"}}
{{- if .ShowAnswer}}{{color "cyan"}} {{ "..." }}{{color "reset"}}{{"\n"}}
{{- else }}
	{{- "  "}}{{- color "cyan"}}[Use arrows to move, space to select, <right> to all, <left> to none, type to filter{{- if and .Help (not .ShowHelp)}}, {{ .Config.HelpInput }} for more help{{end}}]{{color "reset"}}
  {{- "\n"}}
  {{- range $ix, $option := .PageEntries}}
    {{- template "option" $.IterateOption $ix $option}}
  {{- end}}
{{- end}}`

	prompt := survey.MultiSelect{
		Message:  "Select chapters",
		Options:  options,
		VimMode:  viper.GetBool(config.MiniVimMode),
		PageSize: pageSize,
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
