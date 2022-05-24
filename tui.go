package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
	"log"
	"os"
	"sync"
)

/*
Styles
*/

var (
	commonStyle = lipgloss.NewStyle().Margin(2, 2)
)

/*
Key Map
*/

type keyMap struct {
	state     bubbleState
	Quit      key.Binding
	ForceQuit key.Binding

	Select    key.Binding
	SelectAll key.Binding
	Confirm   key.Binding

	Back key.Binding

	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding

	Top    key.Binding
	Bottom key.Binding

	Help key.Binding
}

func (k keyMap) shortHelpFor(state bubbleState) []key.Binding {
	switch state {
	case searchState:
		return []key.Binding{k.Confirm, k.ForceQuit}
	case loadingState:
		return []key.Binding{k.ForceQuit}
	case mangaState:
		return []key.Binding{k.Select, k.Back}
	case chaptersState:
		return []key.Binding{k.Select, k.SelectAll, k.Confirm, k.Back}
	case confirmPromptState:
		return []key.Binding{k.Confirm, k.Back, k.Quit}
	case downloadingState:
		return []key.Binding{k.ForceQuit}
	case exitPromptState:
		return []key.Binding{k.Back, k.Quit}
	}

	return []key.Binding{k.ForceQuit}
}

func (k keyMap) ShortHelp() []key.Binding {
	return k.shortHelpFor(k.state)
}

/*
Model
*/

func newBubble(initialState bubbleState) bubble {
	keys := keyMap{
		state: initialState,

		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit")),
		ForceQuit: key.NewBinding(
			key.WithKeys("ctrl+c", "ctrl+d")),
		Select: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "select")),
		SelectAll: key.NewBinding(
			key.WithKeys("*", "ctrl+a", "tab"),
			key.WithHelp("tab/*", "select all")),
		Confirm: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm")),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back")),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up")),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down")),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "left")),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "right")),
		Top: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "top")),
		Bottom: key.NewBinding(
			key.WithKeys("G"),
			key.WithHelp("G", "bottom")),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help")),
	}

	input := textinput.New()

	spinner_ := spinner.New()
	spinner_.Spinner = spinner.Dot

	progress_ := progress.New(progress.WithDefaultGradient())

	listKeyMap := list.KeyMap{
		CursorUp:             keys.Up,
		CursorDown:           keys.Down,
		NextPage:             keys.Right,
		PrevPage:             keys.Left,
		GoToStart:            keys.Top,
		GoToEnd:              keys.Bottom,
		Filter:               key.Binding{},
		ClearFilter:          key.Binding{},
		CancelWhileFiltering: key.Binding{},
		AcceptWhileFiltering: key.Binding{},
		ShowFullHelp:         keys.Help,
		CloseFullHelp:        keys.Help,
		Quit:                 keys.Quit,
		ForceQuit:            keys.ForceQuit,
	}

	mangaList := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	mangaList.KeyMap = listKeyMap
	mangaList.AdditionalShortHelpKeys = func() []key.Binding { return keys.shortHelpFor(mangaState) }

	chaptersList := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	chaptersList.KeyMap = listKeyMap
	chaptersList.AdditionalShortHelpKeys = func() []key.Binding { return keys.shortHelpFor(chaptersState) }

	bubble_ := bubble{
		state:        initialState,
		keyMap:       keys,
		input:        input,
		spinner:      spinner_,
		mangaList:    mangaList,
		chaptersList: chaptersList,
		progress:     progress_,
		mangaChan:    make(chan []*URL),
		chaptersChan: make(chan []*URL),
		progressChan: make(chan DownloadProgress),
	}

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 0
		height = 0
	}

	bubble_.resize(width, height)
	bubble_.input.Focus()
	return bubble_
}

type bubbleState int

const (
	searchState bubbleState = iota + 1
	loadingState
	mangaState
	chaptersState
	confirmPromptState
	downloadingState
	exitPromptState
)

type bubble struct {
	state bubbleState

	keyMap keyMap

	input        textinput.Model
	spinner      spinner.Model
	mangaList    list.Model
	chaptersList list.Model
	progress     progress.Model

	mangaChan    chan []*URL
	chaptersChan chan []*URL
	progressChan chan DownloadProgress
}

type listItem struct {
	selected bool
	url      *URL
}

func (l *listItem) Select() {
	l.selected = !l.selected
}
func (l listItem) Title() string {
	if l.selected {
		return "+ " + l.url.Info
	}

	return l.url.Info
}

func (l listItem) Description() string {
	return l.url.Address
}

func (l listItem) FilterValue() string {
	return l.url.Info
}

/*
Bubble Init
*/

func (b bubble) Init() tea.Cmd {
	return nil
}

/*
Bubble Update
*/

func (b *bubble) resize(width int, height int) {
	x, y := commonStyle.GetFrameSize()
	b.mangaList.SetSize(width-x, height-y)
	b.chaptersList.SetSize(width-x, height-y)
}

func (b *bubble) setState(state bubbleState) {
	b.state = state
	b.keyMap.state = state
}

type mangaSearchDoneMsg []*URL

func (b bubble) initMangaSearch(query string) tea.Cmd {
	return func() tea.Msg {
		var (
			manga []*URL
			wg    sync.WaitGroup
		)

		wg.Add(len(UserConfig.Scrapers))

		for _, scraper := range UserConfig.Scrapers {
			go func(s *Scraper) {
				defer wg.Done()

				m, err := s.SearchManga(query)

				if err == nil {
					manga = append(manga, m...)
				}
			}(scraper)
		}

		wg.Wait()
		b.mangaChan <- manga

		return nil
	}
}

func (b bubble) waitForMangaSearchCompletion() tea.Cmd {
	return func() tea.Msg {
		return mangaSearchDoneMsg(<-b.mangaChan)
	}
}

type chapterGetDoneMsg []*URL

func (b bubble) initChaptersGet(manga *URL) tea.Cmd {
	return func() tea.Msg {
		chapters, err := manga.Scraper.GetChapters(manga)

		// TODO: Handle it properly
		if err != nil {
			log.Fatalf("Can't get chapters for %s\n", manga.Info)
		}

		b.chaptersChan <- chapters
		return nil
	}
}

func (b bubble) waitForChaptersGetCompletion() tea.Cmd {
	return func() tea.Msg {
		return chapterGetDoneMsg(<-b.chaptersChan)
	}
}

func (b bubble) handleSearchState(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keyMap.Back):
			return b, tea.Quit
		case key.Matches(msg, b.keyMap.Confirm):
			b.setState(loadingState)

			return b, tea.Batch(b.spinner.Tick, b.initMangaSearch(b.input.Value()), b.waitForMangaSearchCompletion())
		}
	}

	b.input, cmd = b.input.Update(msg)
	return b, cmd
}

func (b bubble) handleLoadingState(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case mangaSearchDoneMsg:
		b.setState(mangaState)

		var items []list.Item
		for _, url := range msg {
			items = append(items, listItem{url: url})
		}
		cmd = b.mangaList.SetItems(items)
		return b, cmd
	}

	b.spinner, cmd = b.spinner.Update(msg)
	return b, cmd
}

func (b bubble) handleMangaState(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keyMap.Quit):
			return b, tea.Quit
		case key.Matches(msg, b.keyMap.Back):
			b.setState(searchState)
			return b, nil
		case key.Matches(msg, b.keyMap.Confirm), key.Matches(msg, b.keyMap.Select):
			selected := b.mangaList.SelectedItem().(listItem)
			cmd = b.mangaList.StartSpinner()

			return b, tea.Batch(cmd, b.initChaptersGet(selected.url), b.waitForChaptersGetCompletion())
		}
	case chapterGetDoneMsg:
		b.setState(chaptersState)

		var items []list.Item

		for _, url := range msg {
			items = append(items, listItem{url: url})
		}

		cmd = b.chaptersList.SetItems(items)
		b.mangaList.StopSpinner()
		return b, cmd
	}

	b.mangaList, cmd = b.mangaList.Update(msg)
	return b, cmd
}

func (b bubble) handleChaptersState(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keyMap.Quit):
			return b, tea.Quit
		case key.Matches(msg, b.keyMap.Back):
			b.setState(mangaState)
			return b, nil
		case key.Matches(msg, b.keyMap.Select):
			// TODO: track selected items
			item := b.chaptersList.SelectedItem().(listItem)
			index := b.chaptersList.Index()
			item.Select()
			cmd = b.chaptersList.SetItem(index, item)
			return b, cmd
		case key.Matches(msg, b.keyMap.SelectAll):
			items := b.chaptersList.Items()
			cmds := make([]tea.Cmd, len(items))

			for i, item := range items {
				it := item.(listItem)
				it.Select()
				cmds[i] = b.chaptersList.SetItem(i, it)
			}

			return b, tea.Batch(cmds...)
		}
	}

	b.chaptersList, cmd = b.chaptersList.Update(msg)
	return b, cmd
}

func (b bubble) handleConfirmPromptState(msg tea.Msg) (tea.Model, tea.Cmd) {
	return b, nil
}

func (b bubble) handleDownloadingState(msg tea.Msg) (tea.Model, tea.Cmd) {
	return b, nil
}

func (b bubble) handleExitPromptState(msg tea.Msg) (tea.Model, tea.Cmd) {
	return b, nil
}

func (b bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.resize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keyMap.ForceQuit):
			return b, tea.Quit
		}
	}

	switch b.state {
	case searchState:
		return b.handleSearchState(msg)
	case loadingState:
		return b.handleLoadingState(msg)
	case mangaState:
		return b.handleMangaState(msg)
	case chaptersState:
		return b.handleChaptersState(msg)
	case confirmPromptState:
		return b.handleConfirmPromptState(msg)
	case downloadingState:
		return b.handleDownloadingState(msg)
	case exitPromptState:
		return b.handleExitPromptState(msg)
	}

	log.Fatal("Unknown state encountered")

	// Unreachable
	return b, nil
}

/*
Bubble Render
*/

func (b bubble) View() string {
	var view string

	switch b.state {
	case searchState:
		view = b.input.View()
	case loadingState:
		view = b.spinner.View()
	case mangaState:
		view = b.mangaList.View()
	case chaptersState:
		view = b.chaptersList.View()
	case confirmPromptState:
		view = ""
	case downloadingState:
		view = b.progress.View()
	case exitPromptState:
		view = ""
	}

	return commonStyle.Render(view)
}
