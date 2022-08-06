package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/metafates/mangal/history"
	"github.com/metafates/mangal/provider"
	"github.com/metafates/mangal/source"
	"github.com/samber/lo"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/exp/slices"
	"path/filepath"
	"time"
)

func (b *statefulBubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case error:
		b.plot = randomPlot()
		b.newState(errorState)
	case tea.WindowSizeMsg:
		b.resize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keymap.forceQuit):
			return b, tea.Quit
		case key.Matches(msg, b.keymap.back):

			switch b.state {
			case searchState:
				b.inputC.SetValue("")
			case chaptersState:
				if b.chaptersC.FilterState() != list.Unfiltered {
					b.chaptersC, cmd = b.chaptersC.Update(msg)
					return b, cmd
				}

				b.chaptersC.ResetSelected()
				b.chaptersC.ResetFilter()
				b.selectedChapters = make(map[*source.Chapter]struct{})
			case mangasState:
				if b.mangasC.FilterState() != list.Unfiltered {
					b.mangasC, cmd = b.mangasC.Update(msg)
					return b, cmd
				}

				b.mangasC.ResetSelected()
				b.mangasC.ResetFilter()
			case historyState:
				if b.historyC.FilterState() != list.Unfiltered {
					b.historyC, cmd = b.historyC.Update(msg)
					return b, cmd
				}
			case sourcesState:
				if b.sourcesC.FilterState() != list.Unfiltered {
					b.sourcesC, cmd = b.sourcesC.Update(msg)
					return b, cmd
				}
			}

			b.previousState()
			b.stopLoading()
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
	}

	b.spinnerC, cmd = b.spinnerC.Update(msg)
	return b, cmd
}

func (b *statefulBubble) updateHistory(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case b.sourcesC.FilterState() == list.Filtering:
			break
		case key.Matches(msg, b.keymap.openURL):
			if b.historyC.SelectedItem() != nil {
				chapter := b.historyC.SelectedItem().(*listItem).internal.(*history.SavedChapter)
				err := open.Run(chapter.URL)
				if err != nil {
					b.newState(errorState)
				}
			}
		case key.Matches(msg, b.keymap.selectOne, b.keymap.confirm):
			selected := b.historyC.SelectedItem().(*listItem).internal.(*history.SavedChapter)
			providers := lo.Map(b.sourcesC.Items(), func(i list.Item, _ int) *provider.Provider {
				return i.(*listItem).internal.(*provider.Provider)
			})

			p, ok := lo.Find(providers, func(p *provider.Provider) bool {
				return p.ID == selected.SourceID
			})

			if !ok {
				b.newState(errorState)
				return b, nil
			}

			src, err := p.CreateSource()
			if err != nil {
				b.newState(errorState)
				return b, nil
			}

			b.selectedSource = src

			chapters, err := src.ChaptersOf(&source.Manga{
				Name:     selected.MangaName,
				URL:      selected.MangaURL,
				Index:    0,
				SourceID: selected.SourceID,
			})

			if err != nil {
				b.newState(errorState)
				return b, nil
			}

			_, index, _ := lo.FindIndexOf(chapters, func(c *source.Chapter) bool {
				return c.URL == selected.URL
			})

			items := make([]list.Item, len(chapters))
			for i, c := range chapters {
				items[i] = &listItem{
					internal:    c,
					title:       c.Name,
					description: c.URL,
				}
			}

			cmd = b.chaptersC.SetItems(items)
			b.chaptersC.Select(index)
			b.newState(chaptersState)
			return b, cmd
		}
	}

	b.historyC, cmd = b.historyC.Update(msg)
	return b, cmd
}

func (b *statefulBubble) updateSources(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case b.sourcesC.FilterState() == list.Filtering:
			break
		case key.Matches(msg, b.keymap.confirm, b.keymap.selectOne):
			if b.sourcesC.SelectedItem() == nil {
				break
			}

			s, err := b.sourcesC.SelectedItem().(*listItem).internal.(*provider.Provider).CreateSource()

			if err != nil {
				b.newState(errorState)
			} else {
				b.selectedSource = s
				b.newState(searchState)
			}
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
		case b.mangasC.FilterState() == list.Filtering:
			break
		case key.Matches(msg, b.keymap.confirm, b.keymap.selectOne):
			if b.mangasC.SelectedItem() == nil {
				break
			}

			m, _ := b.mangasC.SelectedItem().(*listItem).internal.(*source.Manga)
			b.selectedManga = m
			return b, tea.Batch(b.getChapters(m), b.waitForChapters(), b.startLoading())
		case key.Matches(msg, b.keymap.openURL):
			if b.mangasC.SelectedItem() == nil {
				break
			}

			m, _ := b.mangasC.SelectedItem().(*listItem).internal.(*source.Manga)
			err := open.Start(m.URL)
			if err != nil {
				b.newState(errorState)
			}
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
	}

	b.mangasC, cmd = b.mangasC.Update(msg)
	return b, cmd
}

func (b *statefulBubble) updateChapters(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case b.chaptersC.FilterState() == list.Filtering:
			break
		case key.Matches(msg, b.keymap.openURL):
			if b.chaptersC.SelectedItem() == nil {
				break
			}

			chapter := b.chaptersC.SelectedItem().(*listItem).internal.(*source.Chapter)
			err := open.Start(chapter.URL)
			if err != nil {
				b.newState(errorState)
			}
		case key.Matches(msg, b.keymap.selectOne):
			if b.chaptersC.SelectedItem() == nil {
				break
			}

			item := b.chaptersC.SelectedItem().(*listItem)
			chapter := item.internal.(*source.Chapter)

			item.toggleMark()
			if item.marked {
				b.selectedChapters[chapter] = struct{}{}
			} else {
				delete(b.selectedChapters, chapter)
			}
		case key.Matches(msg, b.keymap.selectAll):
			items := b.chaptersC.Items()
			if len(items) == 0 {
				break
			}

			for _, item := range items {
				item := item.(*listItem)
				item.marked = true
				chapter := item.internal.(*source.Chapter)
				b.selectedChapters[chapter] = struct{}{}
			}
		case key.Matches(msg, b.keymap.clearSelection):
			items := b.chaptersC.Items()
			if len(items) == 0 {
				break
			}

			for _, item := range items {
				item := item.(*listItem)
				item.marked = false
				chapter := item.internal.(*source.Chapter)
				delete(b.selectedChapters, chapter)
			}
		case key.Matches(msg, b.keymap.read):
			if b.chaptersC.SelectedItem() == nil {
				break
			}

			chapter := b.chaptersC.SelectedItem().(*listItem).internal.(*source.Chapter)
			b.newState(readState)
			return b, tea.Batch(b.readChapter(chapter), b.waitForChapterRead(), b.startLoading())
		case key.Matches(msg, b.keymap.confirm):
			if len(b.selectedChapters) != 0 {
				b.newState(confirmState)
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
		case key.Matches(msg, b.keymap.quit):
			return b, tea.Quit
		case key.Matches(msg, b.keymap.confirm):
			chapters := lo.Keys(b.selectedChapters)
			slices.SortFunc(chapters, func(a, b *source.Chapter) bool {
				return a.Index > b.Index
			})

			for _, chapter := range chapters {
				b.chaptersToDownload.Push(chapter)
			}
			b.newState(downloadState)
			return b, tea.Batch(b.startLoading(), b.downloadChapter(b.chaptersToDownload.Pop()), b.waitForChapterDownload(), b.progressC.SetPercent(0))
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
		inc := 1 / float64(len(b.selectedChapters))

		if b.chaptersToDownload.Length() == 0 {
			// a little hack to make the progress render to the end
			go func() {
				time.Sleep(time.Millisecond * 400)
				b.newState(downloadDoneState)
			}()

			return b, b.progressC.IncrPercent(inc)
		}

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
		case key.Matches(msg, b.keymap.quit):
			return b, tea.Quit
		case key.Matches(msg, b.keymap.openFolder):
			err := open.Start(filepath.Dir(b.lastDownloadedChapterPath))
			if err != nil {
				b.newState(errorState)
			}
		}
	}

	return b, cmd
}

func (b *statefulBubble) updateError(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keymap.quit):
			return b, tea.Quit
		}
	}

	return b, cmd
}
