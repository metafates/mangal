package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/metafates/mangal/icon"
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
	selectedChapters map[*source.Chapter]struct{} // mathematical set

	foundMangasChannel     chan []*source.Manga
	foundChaptersChannel   chan []*source.Chapter
	chapterReadChannel     chan struct{}
	chapterDownloadChannel chan struct{}
	errorChannel           chan error

	progressStatus string

	chaptersToDownload util.Stack[*source.Chapter]

	lastDownloadedChapterPath string
}

func (b *statefulBubble) setState(s state) {
	b.state = s
	b.keymap.setState(s)
}

func (b *statefulBubble) newState(s state) {
	// Transitioning to these states is not allowed (it makes no sense)
	if !lo.Contains([]state{loadingState, exitState, readState, downloadDoneState, downloadState, exitState, confirmState}, b.state) {
		b.statesHistory.Push(b.state)
	}

	b.setState(s)
}

func (b *statefulBubble) previousState() {
	if b.statesHistory.Length() > 0 {
		b.setState(b.statesHistory.Pop())
	}
}

func (b *statefulBubble) resize(width, height int) {
	b.historyC.SetSize(width, height)
	b.sourcesC.SetSize(width, height)
	b.mangasC.SetSize(width, height)
	b.chaptersC.SetSize(width, height)
}

func (b *statefulBubble) startLoading() tea.Cmd {
	b.loading = true
	return tea.Batch(b.mangasC.StartSpinner(), b.chaptersC.StartSpinner())
}

func (b *statefulBubble) stopLoading() tea.Cmd {
	b.loading = false
	b.mangasC.StopSpinner()
	b.chaptersC.StopSpinner()
	return nil
}

func newBubble() *statefulBubble {
	keymap := newStatefulKeymap()
	bubble := statefulBubble{
		state:         idle,
		statesHistory: util.Stack[state]{},
		keymap:        keymap,

		foundMangasChannel:     make(chan []*source.Manga),
		foundChaptersChannel:   make(chan []*source.Chapter),
		chapterReadChannel:     make(chan struct{}),
		chapterDownloadChannel: make(chan struct{}),
		errorChannel:           make(chan error),

		selectedChapters:   make(map[*source.Chapter]struct{}),
		chaptersToDownload: util.Stack[*source.Chapter]{},
	}

	defer func() {
	}()

	makeList := func(title string) list.Model {
		listC := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
		listC.KeyMap = bubble.keymap.forList()
		listC.AdditionalShortHelpKeys = bubble.keymap.ShortHelp
		listC.AdditionalFullHelpKeys = func() []key.Binding {
			return bubble.keymap.FullHelp()[0]
		}
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
	customProviders, err := provider.CustomProviders()

	var items []list.Item
	for _, p := range providers {
		items = append(items, &listItem{
			title:       p.Name,
			description: "Built-in provider " + icon.Get(icon.Go),
			internal:    p,
		})
	}

	if err == nil {
		for _, p := range customProviders {
			items = append(items, &listItem{
				title:       p.Name,
				description: "Custom provider " + icon.Get(icon.Lua),
				internal:    p,
			})
		}
	} else {
		log.Println(err)
	}

	return b.sourcesC.SetItems(items)
}
