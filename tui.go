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
	"sort"
	"strconv"
	"strings"
	"sync"
)

/*
Styles
*/

var (
	commonStyle         = lipgloss.NewStyle().Margin(2, 2)
	accentStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	boldStyle           = lipgloss.NewStyle().Bold(true)
	inputPromptStyle    = accentStyle.Copy().Bold(true)
	inputTitleStyle     = inputPromptStyle.Copy()
	successStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))
	failStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	mangaListTitleStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#9f86c0")).
				Foreground(lipgloss.Color("#231942")).
				Padding(0, 1)
	chaptersListTitleStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#e0b1cb")).
				Foreground(lipgloss.Color("#231942")).
				Padding(0, 1)
)

// keyMap is a map of key bindings for the bubble.
type keyMap struct {
	state bubbleState

	Quit      key.Binding
	ForceQuit key.Binding
	Select    key.Binding
	SelectAll key.Binding
	Confirm   key.Binding
	Open      key.Binding
	Read      key.Binding
	Retry     key.Binding
	Back      key.Binding

	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding

	Top    key.Binding
	Bottom key.Binding

	Help key.Binding
}

// shortHelpFor returns a short list of key bindings for the given state.
func (k keyMap) shortHelpFor(state bubbleState) []key.Binding {
	switch state {
	case searchState:
		return []key.Binding{k.Confirm, k.ForceQuit}
	case loadingState:
		return []key.Binding{k.ForceQuit}
	case mangaState:
		return []key.Binding{k.Open, k.Select, k.Back}
	case chaptersState:
		return []key.Binding{k.Open, k.Read, k.Select, k.SelectAll, k.Confirm, k.Back}
	case confirmPromptState:
		return []key.Binding{k.Confirm, k.Back, k.Quit}
	case downloadingState:
		return []key.Binding{k.ForceQuit}
	case exitPromptState:
		k.Open.SetHelp("o", "open folder")
		k.Retry.SetHelp("r", "redownload failed")
		return []key.Binding{k.Back, k.Open, k.Retry, k.Quit}
	}

	return []key.Binding{k.ForceQuit}
}

// fulleHelpFor returns a full list of key bindings for the given state.
func (k keyMap) fullHelpFor(state bubbleState) []key.Binding {
	switch state {
	case searchState:
		return []key.Binding{k.Confirm, k.ForceQuit}
	case loadingState:
		return []key.Binding{k.ForceQuit}
	case mangaState:
		k.Open.SetHelp("o", "open manga url")
		return []key.Binding{k.Open, k.Select, k.Back}
	case chaptersState:
		k.Read.SetHelp("r", fmt.Sprintf("read chapter in the default %s app", string(UserConfig.Format)))
		k.Open.SetHelp("o", "open chapter url")
		k.Confirm.SetHelp("enter", "download selected chapters")
		return []key.Binding{k.Open, k.Read, k.Select, k.SelectAll, k.Confirm, k.Back}
	case confirmPromptState:
		return []key.Binding{k.Confirm, k.Back, k.Quit}
	case downloadingState:
		return []key.Binding{k.ForceQuit}
	case exitPromptState:
		k.Open.SetHelp("o", "open folder")
		k.Retry.SetHelp("r", "redownload failed")
		return []key.Binding{k.Back, k.Open, k.Retry, k.Quit}
	}

	return []key.Binding{k.ForceQuit}
}

// ShortHelp returns a short list of key bindings for the given state.
func (k keyMap) ShortHelp() []key.Binding {
	return k.shortHelpFor(k.state)
}

// FullHelp returns a full list of key bindings for the given state.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.fullHelpFor(k.state)}
}

// NewBubble creates a new bubble.
func NewBubble(initialState bubbleState) Bubble {
	// Create key bindings.
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
		Read: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "read")),
		Retry: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "retry")),
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

	// Create input component
	input := textinput.New()
	input.Placeholder = UserConfig.UI.Placeholder
	input.CharLimit = 50
	input.Prompt = inputPromptStyle.Render(UserConfig.UI.Prompt + " ")

	// Create spinner component
	spinner_ := spinner.New()
	spinner_.Spinner = spinner.Dot
	spinner_.Style = accentStyle

	// Create progress component
	progress_ := progress.New(progress.WithDefaultGradient())

	// keymap for list components
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

	// Create manga list component
	mangaList := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	mangaList.KeyMap = listKeyMap
	mangaList.AdditionalShortHelpKeys = func() []key.Binding { return keys.shortHelpFor(mangaState) }
	mangaList.AdditionalFullHelpKeys = func() []key.Binding { return keys.fullHelpFor(mangaState) }
	mangaList.Styles.Title = mangaListTitleStyle
	mangaList.Styles.Spinner = accentStyle
	mangaList.SetFilteringEnabled(false)

	// Create chapters list component
	chaptersList := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	chaptersList.KeyMap = listKeyMap
	chaptersList.AdditionalShortHelpKeys = func() []key.Binding { return keys.shortHelpFor(chaptersState) }
	chaptersList.AdditionalFullHelpKeys = func() []key.Binding { return keys.fullHelpFor(chaptersState) }
	chaptersList.Styles.Title = chaptersListTitleStyle
	chaptersList.SetFilteringEnabled(false)
	chaptersList.StatusMessageLifetime = Forever

	// Create new bubble
	bubble_ := Bubble{
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

	// Set initial terminal size
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

// Bubble is the main component of the application
type Bubble struct {
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

// listItem is a list item used in the manga and chapters lists
// It contains the URL of the manga/chapter and the title of the manga/chapter
type listItem struct {
	selected bool
	url      *URL
}

func (l *listItem) Select() {
	l.selected = !l.selected
}
func (l listItem) Title() string {
	var title string

	if l.selected {
		title = accentStyle.Bold(true).Render(UserConfig.UI.Mark) + " " + l.url.Info
	} else {
		title = l.url.Info
	}

	// if user set enumeration to false or if it's a manga
	if !UserConfig.UI.EnumerateChapters || l.url.Relation == nil {
		return title
	}

	index := l.url.Index

	return fmt.Sprintf("[%d] %s", index, title)
}

func (l listItem) Description() string {
	return l.url.Address
}

func (l listItem) FilterValue() string {
	return l.url.Info
}

// Init initializes the bubble
func (b Bubble) Init() tea.Cmd {
	return textinput.Blink
}

// resize the bubble
func (b *Bubble) resize(width int, height int) {
	// Set size to minimum for non-fullscreen runtime
	if !UserConfig.UI.Fullscreen {
		b.mangaList.SetSize(0, 0)
		b.chaptersList.SetSize(0, 0)
		return
	}

	x, y := commonStyle.GetFrameSize()
	b.mangaList.SetSize(width-x, height-y)
	b.chaptersList.SetSize(width-x, height-y)
}

// setState sets the state of the bubble
func (b *Bubble) setState(state bubbleState) {
	b.state = state
	b.keyMap.state = state
}

type mangaSearchDoneMsg []*URL

// initMangaSearch initializes the manga search
func (b Bubble) initMangaSearch(query string) tea.Cmd {
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

// waitForMangaSearchCompletion waits for the manga search to finish
func (b Bubble) waitForMangaSearchCompletion() tea.Cmd {
	return func() tea.Msg {
		return mangaSearchDoneMsg(<-b.mangaChan)
	}
}

type chapterGetDoneMsg []*URL
type chapterDownloadProgressMsg ChapterDownloadProgress

// initChaptersGet initializes the chapters get
func (b Bubble) initChaptersGet(manga *URL) tea.Cmd {
	return func() tea.Msg {
		chapters, err := manga.Scraper.GetChapters(manga)

		if err != nil {
			// set to empty list if error occured and notify user
			b.chaptersChan <- make([]*URL, 0)
			return b.chaptersList.NewStatusMessage("Error occured while fetching chapters")()
		}

		if UserConfig.Anilist.Enabled {
			// cache result
			UserConfig.Anilist.Client.ToAnilistURL(manga)
		}

		b.chaptersChan <- chapters
		return nil
	}
}

// waitForChaptersGetCompletion waits for the chapters get to finish
func (b Bubble) waitForChaptersGetCompletion() tea.Cmd {
	return func() tea.Msg {
		return chapterGetDoneMsg(<-b.chaptersChan)
	}
}

// waitForChapterDownloadProgress waits for the chapter download progress to finish
func (b Bubble) waitForChapterDownloadProgress() tea.Cmd {
	return func() tea.Msg {
		return chapterDownloadProgressMsg(<-b.chapterPagesProgressChan)
	}
}

type chaptersDownloadProgressMsg ChaptersDownloadProgress

// initChaptersDownload initializes the chapters download
func (b Bubble) initChaptersDownload(chapters []*URL) tea.Cmd {
	return func() tea.Msg {
		var (
			failed    []*URL
			succeeded []string
			total     = len(chapters)
			path      string
			err       error
		)

		sort.Slice(chapters, func(i int, j int) bool {
			return chapters[i].Index < chapters[j].Index
		})

		// Download chapters
		for i, chapter := range chapters {
			b.chaptersProgressChan <- ChaptersDownloadProgress{
				Failed:    failed,
				Succeeded: succeeded,
				Total:     total,
				Proceeded: Max(i-1, 0),
				Current:   chapter,
				Done:      false,
			}

			path, err = DownloadChapter(chapter, b.chapterPagesProgressChan, false)
			if err == nil {
				if UserConfig.Anilist.MarkDownloaded {
					_ = UserConfig.Anilist.Client.MarkChapter(chapter.Relation, chapter.Index)
				}

				// use path instead of the chapter name since it is used to get manga folder later
				succeeded = append(succeeded, path)
			} else {
				failed = append(failed, chapter)
			}
		}

		// If epub file was used, create it
		if EpubFile != nil {
			EpubFile.SetAuthor(chapters[0].Scraper.Source.Base)
			if err := EpubFile.Write(path); err != nil {
				log.Fatal("Error while making epub. Please, try again")
			}

			// Close epub file
			EpubFile = nil
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

// waitForChaptersDownloadProgress waits for the chapters download progress to finish
func (b Bubble) waitForChaptersDownloadProgress() tea.Cmd {
	return func() tea.Msg {
		return chaptersDownloadProgressMsg(<-b.chaptersProgressChan)
	}
}

type chapterDownloadedToReadMsg ChaptersDownloadProgress

// initChapterDownloadedToRead initializes the chapter downloaded to read
func (b Bubble) initChapterDownloadToRead(chapter *URL) tea.Cmd {
	return func() tea.Msg {
		var (
			failed    []*URL
			succeeded []string
		)

		b.chaptersProgressChan <- ChaptersDownloadProgress{
			Current:   chapter,
			Done:      false,
			Failed:    failed,
			Succeeded: succeeded,
			Total:     1,
			Proceeded: 0,
		}

		path, err := DownloadChapter(chapter, b.chapterPagesProgressChan, true)

		if err != nil {
			failed = append(failed, chapter)
		} else {
			// Mark chapter as read
			if UserConfig.Anilist.Enabled {
				go func() {
					_ = UserConfig.Anilist.Client.MarkChapter(chapter.Relation, chapter.Index)
				}()
			}

			succeeded = append(succeeded, path)
		}

		b.chaptersProgressChan <- ChaptersDownloadProgress{
			Current:   nil,
			Done:      true,
			Failed:    failed,
			Succeeded: succeeded,
			Total:     1,
			Proceeded: 1,
		}

		return nil
	}
}

// waitForChapterToReadDownloaded waits for the chapter to read download to finish
func (b Bubble) waitForChapterToReadDownloaded() tea.Cmd {
	return func() tea.Msg {
		return chapterDownloadedToReadMsg(<-b.chaptersProgressChan)
	}
}

// handleSearchState handles the search state
func (b Bubble) handleSearchState(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keyMap.Back):
			return b, tea.Quit
		case key.Matches(msg, b.keyMap.Confirm) && b.input.Value() != "":
			b.setState(loadingState)
			return b, tea.Batch(
				b.spinner.Tick,
				b.initMangaSearch(b.input.Value()),
				b.waitForMangaSearchCompletion(),
			)
		}
	}

	b.input, cmd = b.input.Update(msg)
	return b, cmd
}

// handleLoadingState handles the loading state
func (b Bubble) handleLoadingState(msg tea.Msg) (tea.Model, tea.Cmd) {
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

// handleMangaState handles the manga state
func (b Bubble) handleMangaState(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keyMap.Back):
			b.loading = false
			b.mangaList.StopSpinner()
			b.mangaList.Select(0)
			b.setState(searchState)
			return b, b.chaptersList.NewStatusMessage("")

		case key.Matches(msg, b.keyMap.Quit):
			return b, tea.Quit
		case key.Matches(msg, b.keyMap.Open):
			item, ok := b.mangaList.SelectedItem().(listItem)
			if ok {
				_ = open.Start(item.url.Address)
			}
		case b.loading:
			// Do nothing if the chapters are loading
			return b, nil
		case key.Matches(msg, b.keyMap.Confirm), key.Matches(msg, b.keyMap.Select):
			selected, ok := b.mangaList.SelectedItem().(listItem)
			if ok {
				cmd = b.mangaList.StartSpinner()
				b.loading = true

				return b, tea.Batch(cmd, b.initChaptersGet(selected.url), b.waitForChaptersGetCompletion())
			}
		}
	case chapterGetDoneMsg:
		b.setState(chaptersState)
		b.loading = false
		var anilistManga *AnilistURL

		if len(msg) > 0 {
			manga := msg[0].Relation
			if UserConfig.Anilist.Enabled {
				anilistManga = UserConfig.Anilist.Client.ToAnilistURL(manga)
			}
			b.chaptersList.Title = "Chapters - " + PrettyTrim(manga.Info, 30)
		} else {
			b.chaptersList.Title = "Chapters"
		}

		var items []list.Item

		// Sort according to chapter index, in ascending order
		sort.Slice(msg, func(i, j int) bool {
			return msg[i].Index < msg[j].Index
		})

		for _, url := range msg {
			items = append(items, listItem{url: url})
		}

		var cmds []tea.Cmd

		cmds = append(cmds, b.chaptersList.SetItems(items))
		if anilistManga != nil {
			cmds = append(cmds, b.chaptersList.NewStatusMessage("AL: "+PrettyTrim(anilistManga.Title, 25)))
		}
		b.mangaList.StopSpinner()
		return b, cmd
	}

	b.mangaList, cmd = b.mangaList.Update(msg)
	return b, cmd
}

// handleChaptersState handles the chapters state
func (b Bubble) handleChaptersState(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return b, cmd
		case key.Matches(msg, b.keyMap.Open):
			item, ok := b.chaptersList.SelectedItem().(listItem)
			if ok {
				_ = open.Start(item.url.Address)
			}
		case key.Matches(msg, b.keyMap.Read):
			chapter, ok := b.chaptersList.SelectedItem().(listItem)

			if ok {
				b.setState(downloadingState)

				return b, tea.Batch(
					b.progress.SetPercent(0),
					b.spinner.Tick,
					b.initChapterDownloadToRead(chapter.url),
					b.waitForChapterToReadDownloaded(),
					b.waitForChapterDownloadProgress(),
				)
			}
		case key.Matches(msg, b.keyMap.Confirm) && len(b.selectedChapters) > 0:
			b.setState(confirmPromptState)
			return b, nil
		case key.Matches(msg, b.keyMap.Select):
			item, ok := b.chaptersList.SelectedItem().(listItem)
			if ok {
				index := b.chaptersList.Index()
				item.Select()

				if item.selected {
					b.selectedChapters[index] = nil
				} else {
					delete(b.selectedChapters, index)
				}

				cmds := []tea.Cmd{
					b.chaptersList.SetItem(index, item),
				}

				return b, tea.Batch(cmds...)
			}
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

// handleConfirmPromptState handles the confirmation prompt state
func (b Bubble) handleConfirmPromptState(msg tea.Msg) (tea.Model, tea.Cmd) {
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

			return b, tea.Batch(
				b.progress.SetPercent(0),
				b.spinner.Tick,
				b.initChaptersDownload(chapters),
				b.waitForChaptersDownloadProgress(),
				b.waitForChapterDownloadProgress(),
			)
		}
	}

	return b, nil
}

// handleDownloadingState handles the downloading state
func (b Bubble) handleDownloadingState(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case chapterDownloadedToReadMsg:
		b.chaptersDownloadProgressInfo = ChaptersDownloadProgress(msg)

		if msg.Done {
			b.setState(exitPromptState)

			if len(msg.Succeeded) != 0 {
				toRead := msg.Succeeded[0]

				if UserConfig.UseCustomReader {
					_ = open.StartWith(toRead, UserConfig.CustomReader)
				} else {
					_ = open.Start(toRead)
				}
			}

			return b, nil
		}

		cmd = b.progress.SetPercent(float64(len(msg.Succeeded)) / float64(msg.Total))

		return b, tea.Batch(cmd, b.waitForChapterToReadDownloaded(), b.waitForChapterDownloadProgress())
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

// handleExitPromptState handles the exit prompt state
func (b Bubble) handleExitPromptState(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keyMap.Quit):
			RemoveTemp()
			return b, tea.Quit
		case key.Matches(msg, b.keyMap.Back):
			b.setState(chaptersState)
			return b, nil
		case key.Matches(msg, b.keyMap.Retry):
			failed := b.chaptersDownloadProgressInfo.Failed

			if len(failed) > 0 {
				b.setState(downloadingState)
				return b, tea.Batch(
					b.progress.SetPercent(0),
					b.spinner.Tick,
					b.initChaptersDownload(failed),
					b.waitForChaptersDownloadProgress(),
					b.waitForChapterDownloadProgress(),
				)
			}
		case key.Matches(msg, b.keyMap.Open):
			if paths := b.chaptersDownloadProgressInfo.Succeeded; len(paths) > 0 {
				_ = open.Start(filepath.Dir(paths[0]))
			}
		}
	}

	return b, nil
}

// Update handles the Bubble update
func (b Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.resize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keyMap.ForceQuit):
			RemoveTemp()
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

// View handles how the Bubble should be rendered
func (b Bubble) View() string {
	var view string

	template := viewTemplates[b.state]

	switch b.state {
	case searchState:
		if UserConfig.UI.Title == "" {
			view = b.input.View()
		} else {
			view = fmt.Sprintf(template, inputTitleStyle.Render(UserConfig.UI.Title), b.input.View())
		}
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
		view = fmt.Sprintf(
			template,
			accentStyle.Render(strconv.Itoa(chaptersToDownload)),
			Plural("chapter", chaptersToDownload),
			accentStyle.Render(PrettyTrim(mangaName, 40)),
		)
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

	// Do not add help Bubble at these states, since they already have one
	if Contains([]bubbleState{mangaState, chaptersState}, b.state) {
		return commonStyle.Render(view)
	}

	// Add help view
	return commonStyle.Render(fmt.Sprintf("%s\n\n%s", view, b.help.View(b.keyMap)))
}

// viewTemplates is a map of the templates for the different states
var viewTemplates = map[bubbleState]string{
	searchState:        "%s\n\n%s",
	loadingState:       "%s Searching...",
	mangaState:         "%s",
	chaptersState:      "%s",
	confirmPromptState: "Download %s %s of %s ?",
	downloadingState:   "%s\n\n%s\n\n%s %s",
	exitPromptState:    "Done. %s %s downloaded, %s failed",
}
