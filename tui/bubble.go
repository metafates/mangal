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
	"github.com/metafates/mangal/anilist"
	"github.com/metafates/mangal/color"
	"github.com/metafates/mangal/history"
	"github.com/metafates/mangal/installer"
	key2 "github.com/metafates/mangal/key"
	"github.com/metafates/mangal/provider"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/style"
	"github.com/metafates/mangal/util"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/spf13/viper"
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
	anilistC         list.Model
	progressC        progress.Model
	helpC            help.Model

	selectedProviders map[*provider.Provider]struct{}
	selectedSources   []source.Source
	selectedManga     *source.Manga
	selectedChapters  map[*source.Chapter]struct{} // mathematical set

	scrapersLoadedChannel       chan []*installer.Scraper
	scraperInstalledChannel     chan *installer.Scraper
	sourcesLoadedChannel        chan []source.Source
	foundMangasChannel          chan []*source.Manga
	foundChaptersChannel        chan []*source.Chapter
	fetchedAnilistMangasChannel chan []*anilist.Manga
	closestAnilistMangaChannel  chan *anilist.Manga
	chapterReadChannel          chan struct{}
	chapterDownloadChannel      chan struct{}
	errorChannel                chan error

	progressStatus string

	chaptersToDownload util.Stack[*source.Chapter]

	currentDownloadingChapter *source.Chapter
	lastError                 error

	width, height int
	errorPlot     string

	failedChapters   []*source.Chapter
	succededChapters []*source.Chapter

	searchSuggestion mo.Option[string]
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
		anilistSelectState,
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

	b.anilistC.SetSize(listWidth, listHeight)
	b.anilistC.Help.Width = listWidth

	b.progressC.Width = listWidth

	b.width = styledWidth
	b.height = styledHeight
	b.helpC.Width = listWidth
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

		scrapersLoadedChannel:       make(chan []*installer.Scraper),
		scraperInstalledChannel:     make(chan *installer.Scraper),
		sourcesLoadedChannel:        make(chan []source.Source),
		foundMangasChannel:          make(chan []*source.Manga),
		foundChaptersChannel:        make(chan []*source.Chapter),
		fetchedAnilistMangasChannel: make(chan []*anilist.Manga),
		closestAnilistMangaChannel:  make(chan *anilist.Manga),
		chapterReadChannel:          make(chan struct{}),
		chapterDownloadChannel:      make(chan struct{}),
		errorChannel:                make(chan error),

		selectedProviders:  make(map[*provider.Provider]struct{}),
		selectedChapters:   make(map[*source.Chapter]struct{}),
		chaptersToDownload: util.Stack[*source.Chapter]{},

		failedChapters:   make([]*source.Chapter, 0),
		succededChapters: make([]*source.Chapter, 0),
	}

	type listOptions struct {
		TitleStyle mo.Option[lipgloss.Style]
	}

	makeList := func(title string, description bool, options *listOptions) list.Model {
		delegate := list.NewDefaultDelegate()
		delegate.SetSpacing(viper.GetInt(key2.TUIItemSpacing))
		delegate.ShowDescription = description
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
		if titleStyle, ok := options.TitleStyle.Get(); ok {
			listC.Styles.Title = titleStyle
		}

		//listC.StatusMessageLifetime = time.Second * 5
		listC.StatusMessageLifetime = time.Hour * 999 // forever

		return listC
	}

	bubble.helpC = help.New()

	bubble.spinnerC = spinner.New()
	bubble.spinnerC.Spinner = spinner.Dot
	bubble.spinnerC.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))

	bubble.inputC = textinput.New()
	bubble.inputC.Placeholder = "Search"
	bubble.inputC.CharLimit = 60
	bubble.inputC.Prompt = viper.GetString(key2.TUISearchPromptString)

	bubble.progressC = progress.New(progress.WithDefaultGradient())

	bubble.scrapersInstallC = makeList("Install Scrapers", true, &listOptions{
		TitleStyle: mo.Some(
			style.NewColored("#212529", "#ced4da").Padding(0, 1),
		),
	})
	bubble.scrapersInstallC.SetStatusBarItemName("scraper", "scrapers")

	bubble.historyC = makeList("History", true, &listOptions{})
	bubble.sourcesC.SetStatusBarItemName("chapter", "chapters")

	bubble.sourcesC = makeList("Select Source", true, &listOptions{
		TitleStyle: mo.Some(
			style.NewColored("#fefae0", "#bc6c25").Padding(0, 1),
		),
	})
	bubble.sourcesC.SetStatusBarItemName("source", "sources")

	showURLs := viper.GetBool(key2.TUIShowURLs)
	bubble.mangasC = makeList("Mangas", showURLs, &listOptions{
		TitleStyle: mo.Some(
			style.NewColored("#f2e8cf", "#386641").Padding(0, 1),
		),
	})
	bubble.mangasC.SetStatusBarItemName("manga", "mangas")

	bubble.chaptersC = makeList("Chapters", showURLs, &listOptions{
		TitleStyle: mo.Some(
			style.NewColored("#000814", color.Orange).Padding(0, 1),
		),
	})
	bubble.chaptersC.SetStatusBarItemName("chapter", "chapters")

	bubble.anilistC = makeList("Anilist Mangas", showURLs, &listOptions{
		TitleStyle: mo.Some(
			style.NewColored("#bcbedc", "#2b2d42").Padding(0, 1),
		),
	})
	bubble.anilistC.SetStatusBarItemName("manga", "mangas")

	if w, h, err := util.TerminalSize(); err == nil {
		bubble.resize(w, h)
	}

	bubble.inputC.Focus()

	return &bubble
}

func (b *statefulBubble) loadProviders() tea.Cmd {
	providers := provider.Builtins()
	customProviders := provider.Customs()

	var items []list.Item
	for _, p := range providers {
		items = append(items, &listItem{
			internal: p,
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
			internal: p,
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

	chapters := lo.Values(saved)
	slices.SortFunc(chapters, func(a, b *history.SavedChapter) bool {
		if a.MangaName == b.MangaName {
			return a.Name < b.Name
		}
		return a.MangaName < b.MangaName
	})

	var items []list.Item
	for _, c := range chapters {
		items = append(items, &listItem{
			internal: c,
		})
	}

	return tea.Batch(b.historyC.SetItems(items), b.loadProviders()), nil
}
