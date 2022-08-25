package tui

import (
	"fmt"
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
	"github.com/metafates/mangal/installer"
	"github.com/metafates/mangal/provider"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/util"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
	"strings"
	"time"
)

type statefulBubble struct {
	state         state
	statesHistory util.Stack[state]
	loading       bool

	keymap *statefulKeymap

	// components
	spinnerC         spinner.Model
	inputC           textinput.Model
	scrapersInstallC list.Model
	historyC         list.Model
	sourcesC         list.Model
	mangasC          list.Model
	chaptersC        list.Model
	progressC        progress.Model
	helpC            help.Model

	selectedSource   source.Source
	selectedManga    *source.Manga
	selectedChapters map[*source.Chapter]struct{} // mathematical set

	scrapersLoadedChannel   chan []*installer.Scraper
	scraperInstalledChannel chan *installer.Scraper
	sourceLoadedChannel     chan source.Source
	foundMangasChannel      chan []*source.Manga
	foundChaptersChannel    chan []*source.Chapter
	chapterReadChannel      chan struct{}
	chapterDownloadChannel  chan struct{}
	errorChannel            chan error

	progressStatus string

	chaptersToDownload util.Stack[*source.Chapter]

	currentDownloadingChapter *source.Chapter
	lastDownloadedChapterPath string
	lastError                 error

	width, height int
	errorPlot     string

	failedChapters   []*source.Chapter
	succededChapters []*source.Chapter
}

func (b *statefulBubble) raiseError(err error) {
	b.lastError = err
	b.errorPlot = randomPlot()
	b.newState(errorState)
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
		s := b.statesHistory.Pop()
		if s != 0 {
			b.setState(s)
		}
	}
}

func (b *statefulBubble) resize(width, height int) {
	x, y := paddingStyle.GetFrameSize()
	xx, yy := listExtraPaddingStyle.GetFrameSize()

	styledWidth := width - x
	styledHeight := height - y

	listWidth := width - xx
	listHeight := height - yy

	b.scrapersInstallC.SetSize(listWidth, listHeight)
	b.scrapersInstallC.Help.Width = listWidth

	b.historyC.SetSize(listWidth, listHeight)
	b.historyC.Help.Width = listWidth

	b.sourcesC.SetSize(listWidth, listHeight)
	b.sourcesC.Help.Width = listWidth

	b.mangasC.SetSize(listWidth, listHeight)
	b.mangasC.Help.Width = listWidth

	b.chaptersC.SetSize(listWidth, listHeight)
	b.chaptersC.Help.Width = listWidth

	b.progressC.Width = listWidth

	b.width = styledWidth
	b.height = styledHeight
	b.helpC.Width = styledWidth
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
		statesHistory: util.Stack[state]{},
		keymap:        keymap,

		scrapersLoadedChannel:   make(chan []*installer.Scraper),
		scraperInstalledChannel: make(chan *installer.Scraper),
		sourceLoadedChannel:     make(chan source.Source),
		foundMangasChannel:      make(chan []*source.Manga),
		foundChaptersChannel:    make(chan []*source.Chapter),
		chapterReadChannel:      make(chan struct{}),
		chapterDownloadChannel:  make(chan struct{}),
		errorChannel:            make(chan error),

		selectedChapters:   make(map[*source.Chapter]struct{}),
		chaptersToDownload: util.Stack[*source.Chapter]{},

		failedChapters:   make([]*source.Chapter, 0),
		succededChapters: make([]*source.Chapter, 0),
	}

	defer func() {
	}()

	makeList := func(title string) list.Model {
		delegate := list.NewDefaultDelegate()
		delegate.Styles.SelectedTitle = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder(), false, false, false, true).
			BorderForeground(lipgloss.Color("5")).
			Foreground(lipgloss.Color("5")).
			Padding(0, 0, 0, 1)
		delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.Copy().Foreground(lipgloss.Color("7"))

		delegate.Styles.SelectedDesc = delegate.Styles.SelectedTitle.Copy()

		listC := list.New([]list.Item{}, delegate, 0, 0)
		listC.KeyMap = bubble.keymap.forList()
		listC.AdditionalShortHelpKeys = bubble.keymap.ShortHelp
		listC.AdditionalFullHelpKeys = func() []key.Binding {
			return bubble.keymap.FullHelp()[0]
		}
		listC.Title = title
		listC.Styles.NoItems = paddingStyle
		listC.StatusMessageLifetime = time.Second * 5

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

	bubble.scrapersInstallC = makeList("Install Scrapers")
	bubble.scrapersInstallC.SetStatusBarItemName("scraper", "scrapers")

	bubble.historyC = makeList("History")
	bubble.sourcesC.SetStatusBarItemName("chapter", "chapters")

	bubble.sourcesC = makeList("Select Source")
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
	customProviders := provider.CustomProviders()

	var items []list.Item
	for _, p := range providers {
		items = append(items, &listItem{
			title:       p.Name,
			description: "Built-in provider " + icon.Get(icon.Go),
			internal:    p,
		})
	}
	slices.SortFunc(items, func(a, b list.Item) bool {
		// temporary workaround for placing mangadex second because it is not stable for now
		// but, you know, there is nothing more permanent than a temporary solution
		return strings.Compare(a.FilterValue(), b.FilterValue()) > 0
	})

	var customItems []list.Item
	for _, p := range customProviders {
		customItems = append(customItems, &listItem{
			title:       p.Name,
			description: "Custom provider " + icon.Get(icon.Lua),
			internal:    p,
		})
	}
	slices.SortFunc(customItems, func(a, b list.Item) bool {
		return strings.Compare(a.FilterValue(), b.FilterValue()) < 0
	})

	// built-in providers should come first
	return b.sourcesC.SetItems(append(items, customItems...))
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
			description: fmt.Sprintf("%s : %d / %d", s.Name, s.Index, s.MangaChaptersTotal),
			internal:    s,
		})
	}

	return tea.Batch(b.historyC.SetItems(items), b.loadSources()), nil
}
