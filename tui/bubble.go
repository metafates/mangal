package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/metafates/mangal/provider"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/util"
	"github.com/samber/lo"
	"log"
)

type statefulBubble struct {
	state         state
	statesHistory util.Stack[state]
	loading       bool

	keymap *statefulKeymap

	// components
	spinnerC  spinner.Model
	inputC    textinput.Model
	historyC  list.Model
	sourcesC  list.Model
	mangasC   list.Model
	chaptersC list.Model
	progressC progress.Model
	helpC     help.Model

	selectedSource   source.Source
	selectedManga    *source.Manga
	selectedChapter  *source.Chapter
	selectedChapters map[*source.Chapter]struct{} // mathematical set
}

func (b *statefulBubble) newState(s state) {
	// Transitioning to these states is not allowed (it makes no sense)
	// Ignore idle because idle is the state that is set only when the bubble is created
	if b.state != idle && !lo.Contains([]state{loadingState, exitState, readDownloadState, downloadDoneState}, s) {
		b.statesHistory.Push(&b.state)
	}

	b.state = s
	b.keymap.setState(s)
}

func (b *statefulBubble) previousState() {
	if b.statesHistory.Length() > 0 {
		b.statesHistory.Pop()
		b.state = *b.statesHistory.Peek()
		b.keymap.setState(b.state)
	}
}

func (b *statefulBubble) resize(width, height int) {
	b.historyC.SetSize(width, height)
	b.sourcesC.SetSize(width, height)
	b.mangasC.SetSize(width, height)
	b.chaptersC.SetSize(width, height)
}

func (b *statefulBubble) startLoading() {
	b.loading = true
}

func (b *statefulBubble) stopLoading() {
	b.loading = false
}

func newBubble() *statefulBubble {
	keymap := newStatefulKeymap()
	bubble := statefulBubble{
		state:         idle,
		statesHistory: util.Stack[state]{},
		keymap:        keymap,
	}

	defer func() {
	}()

	makeList := func(title string) list.Model {
		listC := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
		listC.KeyMap = bubble.keymap.forList()
		listC.AdditionalShortHelpKeys = bubble.keymap.shortHelp
		listC.AdditionalFullHelpKeys = bubble.keymap.fullHelp
		listC.Title = title

		return listC
	}

	bubble.helpC = help.New()

	bubble.spinnerC = spinner.New()
	bubble.spinnerC.Spinner = spinner.Dot
	bubble.spinnerC.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))

	bubble.inputC = textinput.New()
	bubble.inputC.Placeholder = "Search"
	bubble.inputC.CharLimit = 60
	bubble.inputC.Prompt = "> "

	bubble.progressC = progress.New(progress.WithDefaultGradient())

	bubble.historyC = makeList("History")

	bubble.sourcesC = makeList("Sources")

	bubble.mangasC = makeList("Mangas")

	bubble.chaptersC = makeList("Chapters")

	if w, h, err := util.TerminalSize(); err == nil {
		bubble.resize(w, h)
	}

	bubble.inputC.Focus()

	return &bubble
}

func (b *statefulBubble) loadSources() tea.Cmd {
	providers := provider.DefaultProviders()
	sources, err := source.AvailableCustomSources()

	var items []list.Item
	for _, p := range providers {
		items = append(items, &listItem{
			title:       p.Name,
			description: "Built-in provider",
			internal:    p,
		})
	}

	if err == nil {
		for name, path := range sources {
			items = append(items, &listItem{
				title:       name,
				description: path,
				internal:    nil,
			})
		}
	} else {
		log.Println(err)
	}

	return b.sourcesC.SetItems(items)
}
