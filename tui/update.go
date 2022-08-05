package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/metafates/mangal/provider"
	"github.com/metafates/mangal/source"
	"github.com/skratchdot/open-golang/open"
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
			b.inputC.SetValue("")

			if b.state == chaptersState {
				b.chaptersC.ResetSelected()
				b.selectedChapters = make(map[*source.Chapter]struct{})
			}

			if b.state == mangasState {
				b.mangasC.ResetSelected()
			}

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
	case readState:
		return b.updateRead(msg)
	case downloadState:
		return b.updateDownload(msg)
	case downloadDoneState:
		return b.updateDownloadDone(msg)
	case exitState:
		return b.updateExit(msg)
	case errorState:
		return b.updateError(msg)
	}

	panic("unreachable")
}

func (b *statefulBubble) updateIdle(msg tea.Msg) (tea.Model, tea.Cmd) {
	panic("idle state must not be reached")
	return b, nil
}

func (b *statefulBubble) updateLoading(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keymap.back):
			b.previousState()
			return b, nil
		}
	case []*source.Manga:
		items := make([]list.Item, len(msg))
		for i, m := range msg {
			items[i] = &listItem{
				internal:    m,
				title:       m.Name,
				description: m.URL,
			}
		}

		cmd = b.mangasC.SetItems(items)
		b.newState(mangasState)
		b.stopLoading()
		return b, cmd
	case error:
		b.newState(errorState)
		return b, nil
	}

	b.spinnerC, cmd = b.spinnerC.Update(msg)
	return b, cmd
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
			s, err := b.sourcesC.SelectedItem().(*listItem).internal.(*provider.Provider).CreateSource()

			if err != nil {
				b.newState(errorState)
			} else {
				b.selectedSource = s
				b.newState(searchState)
			}

			return b, nil
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
			b.startLoading()
			b.newState(loadingState)
			return b, tea.Batch(b.searchManga(b.inputC.Value()), b.waitForMangas(), b.spinnerC.Tick)
		}

	}

	b.inputC, cmd = b.inputC.Update(msg)
	return b, cmd
}

func (b *statefulBubble) updateMangas(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keymap.confirm, b.keymap.selectOne):
			m, _ := b.mangasC.SelectedItem().(*listItem).internal.(*source.Manga)
			b.selectedManga = m
			return b, tea.Batch(b.getChapters(m), b.waitForChapters(), b.startLoading())
		case key.Matches(msg, b.keymap.openURL):
			m, _ := b.mangasC.SelectedItem().(*listItem).internal.(*source.Manga)
			err := open.Start(m.URL)
			if err != nil {
				b.newState(errorState)
			}

			return b, nil
		}
	case []*source.Chapter:
		items := make([]list.Item, len(msg))
		for i, c := range msg {
			items[i] = &listItem{
				internal:    c,
				title:       c.Name,
				description: c.URL,
			}
		}

		cmd = b.chaptersC.SetItems(items)
		b.newState(chaptersState)
		b.stopLoading()
		return b, cmd
	case error:
		b.newState(errorState)
		return b, nil
	}

	b.mangasC, cmd = b.mangasC.Update(msg)
	return b, cmd
}

func (b *statefulBubble) updateChapters(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keymap.selectOne):
			item := b.chaptersC.SelectedItem().(*listItem)
			chapter := item.internal.(*source.Chapter)

			item.toggleMark()
			if item.marked {
				b.selectedChapters[chapter] = struct{}{}
			} else {
				delete(b.selectedChapters, chapter)
			}

			return b, nil
		case key.Matches(msg, b.keymap.selectAll):
			items := b.chaptersC.Items()
			for _, item := range items {
				item := item.(*listItem)
				item.toggleMark()
				chapter := item.internal.(*source.Chapter)
				if item.marked {
					b.selectedChapters[chapter] = struct{}{}
				} else {
					delete(b.selectedChapters, chapter)
				}
			}
			return b, nil
		case key.Matches(msg, b.keymap.read):
			chapter := b.chaptersC.SelectedItem().(*listItem).internal.(*source.Chapter)
			b.newState(readState)
			return b, tea.Batch(b.readChapter(chapter), b.waitForChapterRead(), b.startLoading())
		case key.Matches(msg, b.keymap.confirm):
			if len(b.selectedChapters) != 0 {
				b.newState(confirmState)
				return b, nil
			}
		}
	}

	b.chaptersC, cmd = b.chaptersC.Update(msg)
	return b, cmd
}

func (b *statefulBubble) updateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keymap.confirm):
			for chapter := range b.selectedChapters {
				b.chaptersToDownload.Push(chapter)
			}
			b.newState(downloadState)
			return b, tea.Batch(b.startLoading(), b.downloadChapter(b.chaptersToDownload.Pop()), b.waitForChapterDownload())
		case key.Matches(msg, b.keymap.back):
			b.previousState()
		}
	}

	return b, cmd
}

func (b *statefulBubble) updateRead(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.(type) {
	case struct{}:
		b.stopLoading()
		b.previousState()
	}

	b.spinnerC, cmd = b.spinnerC.Update(msg)
	return b, cmd
}

func (b *statefulBubble) updateDownload(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case struct{}:
		if b.chaptersToDownload.Length() == 0 {
			b.newState(downloadDoneState)
			return b, nil
		}

		inc := 1 / float64(len(b.selectedChapters))

		return b, tea.Batch(b.progressC.IncrPercent(inc), b.downloadChapter(b.chaptersToDownload.Pop()), b.waitForChapterDownload())
	case progress.FrameMsg:
		model, cmd := b.progressC.Update(msg)
		b.progressC = model.(progress.Model)
		return b, cmd
	}

	b.spinnerC, cmd = b.spinnerC.Update(msg)
	return b, cmd
}

func (b *statefulBubble) updateDownloadDone(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keymap.back):
			b.previousState()
			b.stopLoading()
			return b, nil
		case key.Matches(msg, b.keymap.quit):
			return b, tea.Quit
		}
	}

	return b, cmd
}

func (b *statefulBubble) updateExit(msg tea.Msg) (tea.Model, tea.Cmd) {
	return b, tea.Quit
}

func (b *statefulBubble) updateError(msg tea.Msg) (tea.Model, tea.Cmd) {
	return b, tea.Quit
}
