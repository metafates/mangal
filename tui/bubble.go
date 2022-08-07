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
	"github.com/metafates/mangal/history"
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

	sourceLoadedChannel    chan source.Source
	foundMangasChannel     chan []*source.Manga
	foundChaptersChannel   chan []*source.Chapter
	chapterReadChannel     chan struct{}
	chapterDownloadChannel chan struct{}
	errorChannel           chan error

	progressStatus string

	chaptersToDownload util.Stack[*source.Chapter]

	currentDownloadingChapter *source.Chapter
	lastDownloadedChapterPath string
	lastError                 error

	width, height int
	plot          string
}

func (b *statefulBubble) setState(s state) {
	b.state = s
	b.keymap.setState(s)
}

func (b *statefulBubble) newState(s state) {
	// do not push state if it is the same as the current state
	if b.state == s {
		return
	}

	// Transitioning to these states is not allowed (it makes no sense)
	if !lo.Contains([]state{
		idle,
		loadingState,
		readState,
		downloadDoneState,
		downloadState,
		confirmState,
	}, b.state) {
		b.statesHistory.Push(b.state)
	}

	b.setState(s)
}

func (b *statefulBubble) previousState() {
	if b.statesHistory.Len() > 0 {
		b.setState(b.statesHistory.Pop())
	}
}

func (b *statefulBubble) resize(width, height int) {
	x, y := paddingStyle.GetFrameSize()
	xx, yy := listExtraPaddingStyle.GetFrameSize()

	styledWidth := width - x
	styledHeight := height - y

	listWidth := width - xx
	listHeight := height - yy

	b.historyC.SetSize(listWidth, listHeight)
	b.sourcesC.SetSize(listWidth, listHeight)
	b.mangasC.SetSize(listWidth, listHeight)
	b.chaptersC.SetSize(listWidth, listHeight)
	b.progressC.Width = styledWidth
	b.width = styledWidth
	b.height = styledHeight
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

		sourceLoadedChannel:    make(chan source.Source),
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
		listC.Styles.NoItems = paddingStyle

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
	bubble.sourcesC.SetStatusBarItemName("source", "sources")

	bubble.mangasC = makeList("Mangas")
	bubble.mangasC.SetStatusBarItemName("manga", "mangas")

	bubble.chaptersC = makeList("Chapters")
	bubble.chaptersC.SetStatusBarItemName("chapter", "chapters")

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

func (b *statefulBubble) loadHistory() (tea.Cmd, error) {
	saved, err := history.Get()
	if err != nil {
		return nil, err
	}

	var items []list.Item
	for _, s := range saved {
		items = append(items, &listItem{
			title:       s.MangaName,
			description: s.Name,
			internal:    s,
		})
	}

	return tea.Batch(b.historyC.SetItems(items), b.loadSources()), nil
}
