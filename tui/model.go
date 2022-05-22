package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/metafates/mangai/api/downloader"
	"github.com/metafates/mangai/api/scraper"
	"github.com/metafates/mangai/config"
	"golang.org/x/term"
	"os"
)

type sessionState int

const (
	searchState sessionState = iota + 1
	spinnerState
	mangaSelectState
	chaptersSelectState
	promptState
	progressState
	exitPromptState
)

type listItem struct {
	url    scraper.URL
	marked bool
}

func (l listItem) Title() string {
	if l.marked {
		return selectedItemStyle.Render("•") + " " + l.url.Info
	}
	return l.url.Info
}

func (l listItem) Description() string {
	return l.url.Address
}

func (l listItem) FilterValue() string {
	return l.url.Info
}

func (l *listItem) mark() {
	l.marked = !l.marked
}

type Bubble struct {
	input    textinput.Model
	spinner  spinner.Model
	manga    list.Model
	chapters list.Model
	prompt   tea.Model
	progress progress.Model
	help     help.Model

	prevManga   string
	prevChapter string
	panelsCount int
	selected    map[int]interface{}

	keys    map[sessionState]keyMap
	bubbles map[sessionState]tea.Model
	config  config.Config

	sub        chan []*scraper.URL
	tick       chan progressInfo
	panelsChan chan downloader.ChapterDownloadInfo

	state      sessionState
	converting bool
}

// resize manga and chapters lists
func (b *Bubble) resize(width, height int) {
	styleW, styleH := commonStyle.GetFrameSize()
	b.manga.SetSize(width-styleW, height-styleH)
	b.chapters.SetSize(width-styleW, height-styleH)
}

func (b Bubble) stateKeyMap() keyMap {
	return b.keys[b.state]
}

// progressInfo used to track manga download progress
type progressInfo struct {
	percent float64
	text    string
}

// New returns an initialized Bubble.
func New() Bubble {
	var k keyMap

	inputModel := textinput.New()
	inputModel.Focus()

	spinnerModel := spinner.New()
	spinnerModel.Spinner = spinner.Dot
	spinnerModel.Style = spinnerStyle

	mangaModel := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	mangaModel.SetFilteringEnabled(false)
	mangaModel.Title = "Manga"
	mangaModel.SetSpinner(spinner.MiniDot)
	k, _ = stateKeyMap[mangaSelectState]
	mangaModel.KeyMap = list.KeyMap{
		CursorUp:             k.Up,
		CursorDown:           k.Down,
		NextPage:             k.Right,
		PrevPage:             k.Left,
		GoToStart:            k.ToStart,
		GoToEnd:              k.ToEnd,
		Filter:               key.Binding{},
		ClearFilter:          key.Binding{},
		CancelWhileFiltering: key.Binding{},
		AcceptWhileFiltering: key.Binding{},
		ShowFullHelp:         k.ShowFullHelp,
		CloseFullHelp:        k.CloseFullHelp,
		Quit:                 k.Quit,
		ForceQuit:            k.Quit,
	}
	mangaModel.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{k.Select, k.Back}
	}

	chaptersModel := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	chaptersModel.SetFilteringEnabled(false)
	chaptersModel.Title = "Chapters"
	chaptersModel.SetSpinner(spinner.MiniDot)
	k, _ = stateKeyMap[chaptersSelectState]
	chaptersModel.KeyMap = list.KeyMap{
		CursorUp:             k.Up,
		CursorDown:           k.Down,
		NextPage:             k.Right,
		PrevPage:             k.Left,
		GoToStart:            k.ToStart,
		GoToEnd:              k.ToEnd,
		Filter:               key.Binding{},
		ClearFilter:          key.Binding{},
		CancelWhileFiltering: key.Binding{},
		AcceptWhileFiltering: key.Binding{},
		ShowFullHelp:         k.ShowFullHelp,
		CloseFullHelp:        k.CloseFullHelp,
		Quit:                 k.Quit,
		ForceQuit:            k.Quit,
	}
	chaptersModel.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{k.Select, k.Confirm, k.Back}
	}

	progressModel := progress.New(progress.WithDefaultGradient())

	bubble := Bubble{
		state:      searchState,
		input:      inputModel,
		spinner:    spinnerModel,
		manga:      mangaModel,
		chapters:   chaptersModel,
		progress:   progressModel,
		help:       help.New(),
		keys:       stateKeyMap,
		config:     config.Get(),
		selected:   map[int]interface{}{},
		sub:        make(chan []*scraper.URL),
		tick:       make(chan progressInfo),
		panelsChan: make(chan downloader.ChapterDownloadInfo),
	}

	termW, termH, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		termW = 0
		termH = 0
	}

	bubble.resize(termW, termH)

	return bubble
}
