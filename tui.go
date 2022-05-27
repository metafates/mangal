package main

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
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/term"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

/*
Styles
*/

var (
	commonStyle           = lipgloss.NewStyle().Margin(2, 2)
	accentStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	selectedItemMarkStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	inputPromptStyle      = accentStyle.Copy().Bold(true)
	inputTitleStyle       = inputPromptStyle.Copy()
	successStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))
	failStyle             = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	mangaListTitleStyle   = lipgloss.NewStyle().
				Background(lipgloss.Color("30")).
				Foreground(lipgloss.Color("230")).
				Padding(0, 1)
	chaptersListTitleStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("62")).
				Foreground(lipgloss.Color("230")).
				Padding(0, 1)
)

/*
Key Map
*/

type keyMap struct {
	state bubbleState

	Quit      key.Binding
	ForceQuit key.Binding
	Select    key.Binding
	SelectAll key.Binding
	Confirm   key.Binding
	Open      key.Binding
	Back      key.Binding

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
		return []key.Binding{k.Open, k.Select, k.Back}
	case chaptersState:
		return []key.Binding{k.Open, k.Select, k.SelectAll, k.Confirm, k.Back}
	case confirmPromptState:
		return []key.Binding{k.Confirm, k.Back, k.Quit}
	case downloadingState:
		return []key.Binding{k.ForceQuit}
	case exitPromptState:
		return []key.Binding{k.Back, k.Open, k.Quit}
	}

	return []key.Binding{k.ForceQuit}
}

func (k keyMap) ShortHelp() []key.Binding {
	return k.shortHelpFor(k.state)
}

func (k keyMap) FullHelp() [][]key.Binding {
	// TODO: add full help
	return [][]key.Binding{k.ShortHelp()}
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
			key.WithKeys("ctrl+c", "ctrl+d"),
			key.WithHelp("ctrl+c", "quit")),
		Select: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "select")),
		SelectAll: key.NewBinding(
			key.WithKeys("*", "ctrl+a", "tab"),
			key.WithHelp("tab", "select all")),
		Confirm: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm")),
		Open: key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "open")),
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
	input.Placeholder = "Search Manga..."
	input.CharLimit = 50
	input.Prompt = inputPromptStyle.Render("> ")

	spinner_ := spinner.New()
	spinner_.Spinner = spinner.Dot
	spinner_.Style = accentStyle

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
	mangaList.Styles.Title = mangaListTitleStyle
	mangaList.Styles.Spinner = accentStyle

	chaptersList := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	chaptersList.KeyMap = listKeyMap
	chaptersList.AdditionalShortHelpKeys = func() []key.Binding { return keys.shortHelpFor(chaptersState) }
	chaptersList.Styles.Title = chaptersListTitleStyle

	bubble_ := bubble{
		state:                        initialState,
		keyMap:                       keys,
		input:                        input,
		spinner:                      spinner_,
		mangaList:                    mangaList,
		chaptersList:                 chaptersList,
		progress:                     progress_,
		help:                         help.New(),
		mangaChan:                    make(chan []*URL),
		chaptersChan:                 make(chan []*URL),
		chaptersProgressChan:         make(chan ChaptersDownloadProgress),
		chapterPagesProgressChan:     make(chan ChapterDownloadProgress),
		selectedChapters:             make(map[int]interface{}),
		chapterDownloadProgressInfo:  ChapterDownloadProgress{},
		chaptersDownloadProgressInfo: ChaptersDownloadProgress{},
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
	state   bubbleState
	loading bool

	keyMap keyMap

	input        textinput.Model
	spinner      spinner.Model
	mangaList    list.Model
	chaptersList list.Model
	progress     progress.Model
	help         help.Model

	mangaChan                chan []*URL
	chaptersChan             chan []*URL
	chaptersProgressChan     chan ChaptersDownloadProgress
	chapterPagesProgressChan chan ChapterDownloadProgress

	chapterDownloadProgressInfo  ChapterDownloadProgress
	chaptersDownloadProgressInfo ChaptersDownloadProgress

	selectedChapters map[int]interface{}
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
		return selectedItemMarkStyle.Render("+") + " " + l.url.Info
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
	return textinput.Blink
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
type chapterDownloadProgressMsg ChapterDownloadProgress

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

func (b bubble) waitForChapterDownloadProgress() tea.Cmd {
	return func() tea.Msg {
		return chapterDownloadProgressMsg(<-b.chapterPagesProgressChan)
	}
}

type chaptersDownloadProgressMsg ChaptersDownloadProgress

func (b bubble) initChaptersDownload(chapters []*URL) tea.Cmd {
	return func() tea.Msg {
		var (
			failed    []*URL
			succeeded []string
			total     = len(chapters)
		)

		for i, chapter := range chapters {
			b.chaptersProgressChan <- ChaptersDownloadProgress{
				Failed:    failed,
				Succeeded: succeeded,
				Total:     total,
				Proceeded: Max(i-1, 0),
				Current:   chapter,
				Done:      false,
			}

			path, err := DownloadChapter(chapter, b.chapterPagesProgressChan)
			if err == nil {
				// use path instead of the chapter name since it is used to get manga folder later
				succeeded = append(succeeded, path)
			} else {
				failed = append(failed, chapter)
			}
		}

		b.chaptersProgressChan <- ChaptersDownloadProgress{
			Failed:    failed,
			Succeeded: succeeded,
			Total:     total,
			Proceeded: total,
			Current:   nil,
			Done:      true,
		}

		return nil
	}
}

func (b bubble) waitForChaptersDownloadProgress() tea.Cmd {
	return func() tea.Msg {
		return chaptersDownloadProgressMsg(<-b.chaptersProgressChan)
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
		b.mangaList.Title = "Manga - " + PrettyTrim(strings.TrimSpace(b.input.Value()), 30)

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
		case key.Matches(msg, b.keyMap.Back):
			b.loading = false
			b.mangaList.StopSpinner()
			b.mangaList.Select(0)
			b.setState(searchState)
			return b, nil
		case key.Matches(msg, b.keyMap.Quit):
			return b, tea.Quit
		case key.Matches(msg, b.keyMap.Open):
			item := b.mangaList.SelectedItem().(listItem)
			_ = open.Start(item.url.Address)
		case b.loading:
			// Do nothing if the chapters are loading
			return b, nil
		case key.Matches(msg, b.keyMap.Confirm), key.Matches(msg, b.keyMap.Select):
			selected := b.mangaList.SelectedItem().(listItem)
			cmd = b.mangaList.StartSpinner()
			b.loading = true

			return b, tea.Batch(cmd, b.initChaptersGet(selected.url), b.waitForChaptersGetCompletion())
		}
	case chapterGetDoneMsg:
		b.setState(chaptersState)
		b.loading = false

		if len(msg) > 0 {
			b.chaptersList.Title = "Chapters - " + PrettyTrim(msg[0].Relation.Info, 30)
		} else {
			b.chaptersList.Title = "Chapters"
		}

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
			// reset selected items
			b.chaptersList.Select(0)
			b.selectedChapters = make(map[int]interface{})

			b.setState(mangaState)
			return b, nil
		case key.Matches(msg, b.keyMap.Open):
			item := b.chaptersList.SelectedItem().(listItem)
			_ = open.Start(item.url.Address)
		case key.Matches(msg, b.keyMap.Confirm):
			if len(b.selectedChapters) > 0 {
				b.setState(confirmPromptState)
				return b, nil
			}

			return b, nil
		case key.Matches(msg, b.keyMap.Select):
			item := b.chaptersList.SelectedItem().(listItem)
			index := b.chaptersList.Index()
			item.Select()

			if item.selected {
				b.selectedChapters[index] = nil
			} else {
				delete(b.selectedChapters, index)
			}

			cmd = b.chaptersList.SetItem(index, item)
			return b, cmd
		case key.Matches(msg, b.keyMap.SelectAll):
			items := b.chaptersList.Items()
			cmds := make([]tea.Cmd, len(items))

			for i, item := range items {
				it := item.(listItem)
				it.Select()

				if it.selected {
					b.selectedChapters[i] = nil
				} else {
					delete(b.selectedChapters, i)
				}

				cmds[i] = b.chaptersList.SetItem(i, it)
			}

			return b, tea.Batch(cmds...)
		}
	}

	b.chaptersList, cmd = b.chaptersList.Update(msg)
	return b, cmd
}

func (b bubble) handleConfirmPromptState(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keyMap.Quit):
			return b, tea.Quit
		case key.Matches(msg, b.keyMap.Back):
			b.setState(chaptersState)
			return b, nil
		case key.Matches(msg, b.keyMap.Confirm):
			b.setState(downloadingState)

			var chapters []*URL

			items := b.chaptersList.Items()

			for index := range b.selectedChapters {
				chapters = append(chapters, items[index].(listItem).url)
			}

			return b, tea.Batch(b.progress.SetPercent(0), b.spinner.Tick, b.initChaptersDownload(chapters), b.waitForChaptersDownloadProgress(), b.waitForChapterDownloadProgress())
		}
	}

	return b, nil
}

func (b bubble) handleDownloadingState(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case chapterDownloadProgressMsg:
		b.spinner, cmd = b.spinner.Update(msg)
		b.chapterDownloadProgressInfo = ChapterDownloadProgress(msg)
		return b, tea.Batch(cmd, b.waitForChapterDownloadProgress(), b.waitForChaptersGetCompletion())
	case chaptersDownloadProgressMsg:
		b.chaptersDownloadProgressInfo = ChaptersDownloadProgress(msg)

		if msg.Done {
			b.setState(exitPromptState)
			return b, nil
		}

		cmd = b.progress.SetPercent(float64(len(msg.Succeeded)) / float64(msg.Total))
		return b, tea.Batch(cmd, b.waitForChaptersDownloadProgress(), b.waitForChapterDownloadProgress())
	case progress.FrameMsg:
		var p tea.Model
		// ???? why progress.Update() returns tea.Model and not progress.Model?
		p, cmd = b.progress.Update(msg)
		b.progress = p.(progress.Model)
		return b, cmd
	}

	b.spinner, cmd = b.spinner.Update(msg)
	return b, cmd
}

func (b bubble) handleExitPromptState(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keyMap.Quit):
			return b, tea.Quit
		case key.Matches(msg, b.keyMap.Back):
			b.setState(chaptersState)
			return b, nil
		case key.Matches(msg, b.keyMap.Open):
			if paths := b.chaptersDownloadProgressInfo.Succeeded; len(paths) > 0 {
				_ = open.Start(filepath.Dir(paths[0]))
			}
		}
	}

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

	template := viewTemplates[b.state]

	switch b.state {
	case searchState:
		view = fmt.Sprintf(template, inputTitleStyle.Render(AppName), b.input.View())
	case loadingState:
		view = fmt.Sprintf(template, b.spinner.View())
	case mangaState:
		view = fmt.Sprintf(template, b.mangaList.View())
	case chaptersState:
		view = fmt.Sprintf(template, b.chaptersList.View())
	case confirmPromptState:
		// Should be unreachable
		if len(b.selectedChapters) == 0 {
			log.Fatal("No chapters selected")
		}

		mangaName := b.chaptersList.Items()[0].(listItem).url.Relation.Info
		chaptersToDownload := len(b.selectedChapters)
		view = fmt.Sprintf(template, accentStyle.Render(strconv.Itoa(chaptersToDownload)), Plural("chapter", chaptersToDownload), accentStyle.Render(PrettyTrim(mangaName, 40)))
	case downloadingState:

		var header string

		// It shouldn't be nil at this stage but it panics TODO: FIX THIS
		if b.chaptersDownloadProgressInfo.Current != nil {
			mangaName := b.chaptersDownloadProgressInfo.Current.Info
			header = fmt.Sprintf("Downloading %s", PrettyTrim(accentStyle.Render(mangaName), 40))
		} else {
			header = "Preparing for download..."
		}

		subheader := b.chapterDownloadProgressInfo.Message
		view = fmt.Sprintf("%s\n\n%s\n\n%s %s", header, b.progress.View(), b.spinner.View(), subheader)
	case exitPromptState:
		succeeded := b.chaptersDownloadProgressInfo.Succeeded
		failed := b.chaptersDownloadProgressInfo.Failed

		succeededRendered := successStyle.Render(strconv.Itoa(len(succeeded)))
		failedRendered := failStyle.Render(strconv.Itoa(len(failed)))

		view = fmt.Sprintf(template, succeededRendered, Plural("chapter", len(succeeded)), failedRendered)

		// show failed chapters
		for _, chapter := range failed {
			view += fmt.Sprintf("\n\n%s %s", failStyle.Render("Failed"), chapter.Info)
		}
	}

	// Do not add help bubble at these states, since they already have one
	if Contains([]bubbleState{mangaState, chaptersState}, b.state) {
		return commonStyle.Render(view)
	}

	// Add help view
	return commonStyle.Render(fmt.Sprintf("%s\n\n%s", view, b.help.View(b.keyMap)))
}

var viewTemplates = map[bubbleState]string{
	searchState:        "%s\n\n%s",
	loadingState:       "%s Searching...",
	mangaState:         "%s",
	chaptersState:      "%s",
	confirmPromptState: "Download %s %s of %s ?",
	downloadingState:   "%s\n\n%s\n\n%s %s",
	exitPromptState:    "Done. %s %s downloaded, %s failed",
}
