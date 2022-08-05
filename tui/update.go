package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (b *statefulBubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.resize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keymap.forceQuit):
			return b, tea.Quit
		case key.Matches(msg, b.keymap.back):
			b.previousState()
			return b, nil
		}
	}

	switch b.state {
	case idle:
		return b.updateIdle(msg)
	case loadingState:
		return b.updateLoading(msg)
	case historyState:
		return b.updateHistory(msg)
	case sourcesState:
		return b.updateSources(msg)
	case searchState:
		return b.updateSearch(msg)
	case mangasState:
		return b.updateMangas(msg)
	case chaptersState:
		return b.updateChapters(msg)
	case confirmState:
		return b.updateConfirm(msg)
	case readDownloadState:
		return b.updateReadDownload(msg)
	case readDownloadDoneState:
		return b.updateReadDownloadDone(msg)
	case downloadState:
		return b.updateDownload(msg)
	case downloadDoneState:
		return b.updateDownloadDone(msg)
	case exitState:
		return b.updateExit(msg)
	}

	panic("unreachable")
}

func (b *statefulBubble) updateIdle(msg tea.Msg) (tea.Model, tea.Cmd) {
	panic("idle state must not be reached")
	return b, nil
}

func (b *statefulBubble) updateLoading(msg tea.Msg) (tea.Model, tea.Cmd) {
	return b, nil
}

func (b *statefulBubble) updateHistory(msg tea.Msg) (tea.Model, tea.Cmd) {
	return b, nil
}

func (b *statefulBubble) updateSources(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keymap.confirm, b.keymap.selectOne):
			return b, tea.Quit
		}
	}

	b.sourcesC, cmd = b.sourcesC.Update(msg)
	return b, cmd
}

func (b *statefulBubble) updateSearch(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keymap.confirm) && b.inputC.Value() != "":
			return b, tea.Quit
		}

	}

	b.inputC, cmd = b.inputC.Update(msg)
	return b, cmd
}

func (b *statefulBubble) updateMangas(msg tea.Msg) (tea.Model, tea.Cmd) {
	return b, nil
}

func (b *statefulBubble) updateChapters(msg tea.Msg) (tea.Model, tea.Cmd) {
	return b, nil
}

func (b *statefulBubble) updateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	return b, nil
}

func (b *statefulBubble) updateReadDownload(msg tea.Msg) (tea.Model, tea.Cmd) {
	return b, nil
}

func (b *statefulBubble) updateReadDownloadDone(msg tea.Msg) (tea.Model, tea.Cmd) {
	return b, nil
}

func (b *statefulBubble) updateDownload(msg tea.Msg) (tea.Model, tea.Cmd) {
	return b, nil
}

func (b *statefulBubble) updateDownloadDone(msg tea.Msg) (tea.Model, tea.Cmd) {
	return b, nil
}

func (b *statefulBubble) updateExit(msg tea.Msg) (tea.Model, tea.Cmd) {
	return b, nil
}
