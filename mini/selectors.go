package mini

import (
	"errors"
	"github.com/AlecAivazis/survey/v2"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/icon"
	"github.com/metafates/mangal/provider"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/style"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

func selectSource() (source.Source, error) {

	defaultProviders := provider.DefaultProviders()
	customProviders, err := provider.CustomProviders()

	if err != nil {
		return nil, err
	}

	var sources = make(map[string]func() (source.Source, error))

	for name, p := range customProviders {
		sources[name+" "+icon.Get(icon.Lua)] = p.CreateSource
	}

	for name, p := range defaultProviders {
		sources[name+" "+icon.Get(icon.Go)] = p.CreateSource
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
	if len(mangas) == 0 {
		return nil, errors.New("no manga found")
	}

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

func selectChapter(chapters []*source.Chapter, offset int) (*source.Chapter, error) {
	if len(chapters) == 0 {
		return nil, errors.New("no chapters found")
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
		Default:  options[offset],
	}

	var chapterName string
	err := survey.AskOne(&prompt, &chapterName, survey.WithValidator(survey.Required))
	if err != nil {
		return nil, err
	}

	return c[chapterName], nil
}

func selectChapters(chapters []*source.Chapter) ([]*source.Chapter, error) {
	if len(chapters) == 0 {
		return nil, errors.New("no chapters found")
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
	err := survey.AskOne(&prompt, &chapterNames, survey.WithValidator(survey.Required))
	if err != nil {
		return nil, err
	}

	var chaptersToDownload = make([]*source.Chapter, len(chapterNames))
	for i, chapterName := range chapterNames {
		chaptersToDownload[i] = c[chapterName]
	}

	slices.SortFunc(chaptersToDownload, func(a, b *source.Chapter) bool {
		return a.Index < b.Index
	})

	return chaptersToDownload, nil
}
