package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/util"
	"github.com/samber/lo"
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

func (b *statefulBubble) setState(newState state) {
	// Transitioning to these states is not allowed (it makes no sense)
	// Ignore idle because idle is the state that is set only when the bubble is created
	if b.state != idle && !lo.Contains([]state{loadingState, exitState, readDownloadState, downloadDoneState}, newState) {
		b.statesHistory.Push(&b.state)
	}

	b.state = newState
	b.keymap.setState(newState)
}

func (b *statefulBubble) resize(width, height int) {
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
		if w, h, err := util.TerminalSize(); err != nil {
			bubble.resize(0, 0)
		} else {
			bubble.resize(w, h)
		}

		bubble.inputC.Focus()
	}()

	makeList := func(title string) list.Model {
		listC := list.New(nil, list.NewDefaultDelegate(), 0, 0)
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
	bubble.inputC.CharLimit = 40
	bubble.inputC.Prompt = "> "

	bubble.progressC = progress.New(progress.WithDefaultGradient())

	bubble.historyC = makeList("History")

	bubble.sourcesC = makeList("Sources")

	bubble.chaptersC = makeList("Chapters")

	return &bubble
}
