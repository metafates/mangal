package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/metafates/mangal/common"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/downloader"
	scraper2 "github.com/metafates/mangal/scraper"
	"github.com/metafates/mangal/style"
	"golang.org/x/term"
	"os"
)

// Bubble is the main component of the application
type Bubble struct {
	state   bubbleState
	loading bool

	keyMap *keyMap

	input        *textinput.Model
	spinner      *spinner.Model
	ResumeList   *list.Model
	mangaList    *list.Model
	chaptersList *list.Model
	progress     *progress.Model
	help         *help.Model

	mangaChan                chan []*scraper2.URL
	chaptersChan             chan []*scraper2.URL
	chaptersProgressChan     chan downloader.ChaptersDownloadProgress
	chapterPagesProgressChan chan downloader.ChapterDownloadProgress

	chapterDownloadProgressInfo  downloader.ChapterDownloadProgress
	chaptersDownloadProgressInfo downloader.ChaptersDownloadProgress

	selectedChapters map[int]interface{}
}

// NewBubble creates a new bubble.
func NewBubble(initialState bubbleState) *Bubble {
	// Create key bindings.
	keys := &keyMap{
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
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter")),
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
	input.Placeholder = config.UserConfig.UI.Placeholder
	input.CharLimit = 50
	input.Prompt = style.InputPromptStyle.Render(config.UserConfig.UI.Prompt + " ")

	// Create spinner component
	spinner_ := spinner.New()
	spinner_.Spinner = spinner.Dot
	spinner_.Style = style.AccentStyle

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
		Filter:               keys.Filter,
		ClearFilter:          key.Binding{},
		CancelWhileFiltering: keys.Back,
		AcceptWhileFiltering: keys.Confirm,
		ShowFullHelp:         keys.Help,
		CloseFullHelp:        keys.Help,
		Quit:                 keys.Quit,
		ForceQuit:            keys.ForceQuit,
	}

	// Create resume list component
	ResumeList := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	ResumeList.KeyMap = listKeyMap
	ResumeList.AdditionalShortHelpKeys = func() []key.Binding { return keys.shortHelpFor(MangaState) }
	ResumeList.AdditionalFullHelpKeys = func() []key.Binding { return keys.fullHelpFor(MangaState) }
	ResumeList.Styles.Title = style.MangaListTitleStyle
	ResumeList.Styles.Spinner = style.AccentStyle
	ResumeList.Title = "Resume"
	ResumeList.SetFilteringEnabled(true)

	// Create manga list component
	mangaList := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	mangaList.KeyMap = listKeyMap
	mangaList.AdditionalShortHelpKeys = func() []key.Binding { return keys.shortHelpFor(MangaState) }
	mangaList.AdditionalFullHelpKeys = func() []key.Binding { return keys.fullHelpFor(MangaState) }
	mangaList.Styles.Title = style.MangaListTitleStyle
	mangaList.Styles.Spinner = style.AccentStyle
	mangaList.SetFilteringEnabled(false)

	// Create chapters list component
	chaptersList := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	chaptersList.KeyMap = listKeyMap
	chaptersList.AdditionalShortHelpKeys = func() []key.Binding { return keys.shortHelpFor(ChaptersState) }
	chaptersList.AdditionalFullHelpKeys = func() []key.Binding { return keys.fullHelpFor(ChaptersState) }
	chaptersList.Styles.Title = style.ChaptersListTitleStyle
	chaptersList.SetFilteringEnabled(false)
	chaptersList.StatusMessageLifetime = common.Forever

	helpModel := help.New()

	// Create new bubble
	bubble_ := Bubble{
		state:                        initialState,
		keyMap:                       keys,
		input:                        &input,
		spinner:                      &spinner_,
		ResumeList:                   &ResumeList,
		mangaList:                    &mangaList,
		chaptersList:                 &chaptersList,
		progress:                     &progress_,
		help:                         &helpModel,
		mangaChan:                    make(chan []*scraper2.URL),
		chaptersChan:                 make(chan []*scraper2.URL),
		chaptersProgressChan:         make(chan downloader.ChaptersDownloadProgress),
		chapterPagesProgressChan:     make(chan downloader.ChapterDownloadProgress),
		selectedChapters:             make(map[int]interface{}),
		chapterDownloadProgressInfo:  downloader.ChapterDownloadProgress{},
		chaptersDownloadProgressInfo: downloader.ChaptersDownloadProgress{},
	}

	// Set initial terminal size
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 0
		height = 0
	}

	bubble_.resize(width, height)
	bubble_.input.Focus()
	return &bubble_
}
