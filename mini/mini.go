package mini

import (
	"errors"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/util"
	"github.com/samber/lo"
	"os"
)

var (
	truncateAt = 100
)

type Options struct {
	Download bool
	Continue bool
}

type mini struct {
	width, height int

	state         state
	statesHistory util.Stack[state]

	download bool

	selectedSource source.Source

	cachedMangas   map[string][]*source.Manga
	cachedChapters map[string][]*source.Chapter
	cachedPages    map[string][]*source.Page

	query            string
	selectedManga    *source.Manga
	selectedChapters []*source.Chapter
}

func newMini() *mini {
	return &mini{
		statesHistory:  util.Stack[state]{},
		cachedMangas:   make(map[string][]*source.Manga),
		cachedChapters: make(map[string][]*source.Chapter),
		cachedPages:    make(map[string][]*source.Page),
	}
}

func (m *mini) previousState() {
	if m.statesHistory.Len() > 0 {
		m.setState(m.statesHistory.Pop())
	}
}

func (m *mini) setState(s state) {
	m.state = s
}

func (m *mini) newState(s state) {
	// do not push state if it is the same as the current state
	if m.state == s {
		return
	}

	// Transitioning to these states is not allowed (it makes no sense)
	if !lo.Contains([]state{}, m.state) {
		m.statesHistory.Push(m.state)
	}

	m.setState(s)
}

func Run(options *Options) error {
	if options.Continue && options.Download {
		return errors.New("cannot download and continue")
	}

	m := newMini()
	m.state = sourceSelectState
	if options.Continue {
		m.state = historySelectState
	}

	m.download = options.Download

	if w, h, err := util.TerminalSize(); err == nil {
		m.width, m.height = w, h
		truncateAt = w
	}

	var err error

	for {
		if m.handleState() != nil {
			return err
		}
	}
}

func (m *mini) handleState() error {
	switch m.state {
	case historySelectState:
		return m.handleHistorySelectState()
	case sourceSelectState:
		return m.handleSourceSelectState()
	case mangasSearchState:
		return m.handleMangaSearchState()
	case mangaSelectState:
		return m.handleMangaSelectState()
	case chapterSelectState:
		return m.handleChapterSelectState()
	case chapterReadState:
		return m.handleChapterReadState()
	case chaptersDownloadState:
		return m.handleChaptersDownloadState()
	case quitState:
		os.Exit(0)
	}

	return nil
}
