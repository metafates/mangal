package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/metafates/mangal/history"
	"github.com/metafates/mangal/installer"
	"github.com/metafates/mangal/open"
	"github.com/metafates/mangal/provider"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/util"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
	"time"
)

func (b *statefulBubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case error:
		b.raiseError(msg)
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
			case scrapersInstallState:
				if b.scrapersInstallC.FilterState() != list.Unfiltered {
					b.scrapersInstallC, cmd = b.scrapersInstallC.Update(msg)
					return b, cmd
				}
			}

			b.previousState()
			b.stopLoading()
			b.failedChapters = make([]*source.Chapter, 0)
			b.succededChapters = make([]*source.Chapter, 0)
			return b, nil
		}
	}

	switch b.state {
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
	case scrapersInstallState:
		return b.updateScrapersInstall(msg)
	case errorState:
		return b.updateError(msg)
	}

	panic("unreachable")
}

func (b *statefulBubble) updateScrapersInstall(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if len(b.scrapersInstallC.Items()) == 0 {
		b.newState(loadingState)
		return b, tea.Batch(b.startLoading(), b.loadScrapers(), b.waitForScrapersLoaded())
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case b.scrapersInstallC.FilterState() == list.Filtering:
			break
		case key.Matches(msg, b.keymap.openURL):
			url := b.scrapersInstallC.SelectedItem().(*listItem).internal.(*installer.Scraper).GithubURL()
			err := open.Run(url)
			if err != nil {
				b.lastError = err
				b.newState(errorState)
			}
		case key.Matches(msg, b.keymap.selectOne, b.keymap.confirm):
			scraper := b.scrapersInstallC.SelectedItem().(*listItem).internal.(*installer.Scraper)
			b.newState(loadingState)
			return b, tea.Batch(b.startLoading(), b.installScraper(scraper), b.waitForScraperInstallation())
		}
	}

	b.scrapersInstallC, cmd = b.scrapersInstallC.Update(msg)
	return b, cmd
}

func (b *statefulBubble) updateLoading(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds = make([]tea.Cmd, 0)
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keymap.back):
			b.previousState()
		}
	case []*installer.Scraper:
		b.newState(scrapersInstallState)
		return b, b.stopLoading()
	case *installer.Scraper:
		b.newState(scrapersInstallState)
		b.scrapersInstallC.NewStatusMessage(fmt.Sprintf("Installed %s", msg.Name))
		return b, b.stopLoading()
	case []*source.Manga:
		items := make([]list.Item, len(msg))
		for i, m := range msg {
			items[i] = &listItem{
				internal:    m,
				title:       m.Name,
				description: m.URL,
			}
		}

		cmds = append(cmds, b.mangasC.SetItems(items))
		b.newState(mangasState)
		b.stopLoading()
	case []*source.Chapter:
		if b.statesHistory.Peek() == historyState {
			b.newState(historyState)
			b.stopLoading()
			cmds = append(cmds, func() tea.Msg {
				return msg
			})
		}
	case []source.Source:
		b.selectedSources = msg

		if b.statesHistory.Peek() == historyState {
			b.newState(historyState)
			b.stopLoading()
			cmds = append(cmds, func() tea.Msg {
				return msg
			})
		} else {
			b.stopLoading()
			b.newState(searchState)
		}
	}

	b.spinnerC, cmd = b.spinnerC.Update(msg)
	return b, tea.Batch(append(cmds, cmd)...)
}

func (b *statefulBubble) updateHistory(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case []source.Source:
		b.selectedSources = msg
		selected := b.historyC.SelectedItem().(*listItem).internal.(*history.SavedChapter)

		manga := &source.Manga{
			Name:   selected.MangaName,
			URL:    selected.MangaURL,
			Index:  0,
			ID:     selected.MangaID,
			Source: b.selectedSources[0],
		}

		b.selectedManga = manga
		b.newState(loadingState)
		return b, tea.Batch(
			b.startLoading(),
			b.getChapters(manga),
			b.waitForChapters(),
		)
	case []*source.Chapter:
		items := make([]list.Item, len(msg))
		selected := b.historyC.SelectedItem().(*listItem).internal.(*history.SavedChapter)

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
		selectCmd := b.selectChapterBy(func(chapter *source.Chapter) bool {
			return chapter.URL == selected.URL
		})
		return b, tea.Batch(cmd, selectCmd)
	case tea.KeyMsg:
		switch {
		case b.historyC.FilterState() == list.Filtering:
			break
		case key.Matches(msg, b.keymap.openURL):
			if b.historyC.SelectedItem() != nil {
				chapter := b.historyC.SelectedItem().(*listItem).internal.(*history.SavedChapter)
				err := open.Run(chapter.URL)
				if err != nil {
					b.raiseError(err)
				}
			}
		case key.Matches(msg, b.keymap.remove):
			if b.historyC.SelectedItem() != nil {
				chapter := b.historyC.SelectedItem().(*listItem).internal.(*history.SavedChapter)
				_ = history.Remove(chapter)
				cmd, err := b.loadHistory()
				if err != nil {
					return nil, nil
				}

				return b, cmd
			}
		case key.Matches(msg, b.keymap.selectOne, b.keymap.confirm):
			if b.historyC.SelectedItem() != nil {
				selected := b.historyC.SelectedItem().(*listItem).internal.(*history.SavedChapter)
				providers := lo.Map(b.sourcesC.Items(), func(i list.Item, _ int) *provider.Provider {
					return i.(*listItem).internal.(*provider.Provider)
				})

				p, ok := lo.Find(providers, func(p *provider.Provider) bool {
					return p.ID == selected.SourceID
				})

				if !ok {
					err := fmt.Errorf("provider %s not found", selected.SourceID)
					b.raiseError(err)
					return b, nil
				}

				b.newState(loadingState)
				return b, tea.Batch(b.startLoading(), b.loadSources([]*provider.Provider{p}), b.waitForSourcesLoaded())
			}
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
		case key.Matches(msg, b.keymap.selectAll):
			for _, item := range b.sourcesC.Items() {
				item := item.(*listItem)
				item.marked = true
				b.selectedProviders[item.internal.(*provider.Provider)] = struct{}{}
			}
		case key.Matches(msg, b.keymap.clearSelection):
			for _, item := range b.sourcesC.Items() {
				item := item.(*listItem)
				item.marked = false
				delete(b.selectedProviders, item.internal.(*provider.Provider))
			}
		case key.Matches(msg, b.keymap.selectOne):
			if b.sourcesC.SelectedItem() == nil {
				break
			}

			item := b.sourcesC.SelectedItem().(*listItem)
			p := item.internal.(*provider.Provider)

			if item.marked {
				delete(b.selectedProviders, p)
			} else {
				b.selectedProviders[p] = struct{}{}
			}

			item.toggleMark()
		case key.Matches(msg, b.keymap.confirm):
			if b.sourcesC.SelectedItem() == nil {
				break
			}

			item := b.sourcesC.SelectedItem().(*listItem)

			if len(b.selectedProviders) == 0 {
				b.selectedProviders[item.internal.(*provider.Provider)] = struct{}{}
			}

			b.newState(loadingState)
			return b, tea.Batch(b.startLoading(), b.loadSources(lo.Keys(b.selectedProviders)), b.waitForSourcesLoaded())
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
				b.raiseError(err)
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
				b.raiseError(err)
			}
		case key.Matches(msg, b.keymap.selectVolume):
			if b.chaptersC.SelectedItem() == nil {
				break
			}

			chapter := b.chaptersC.SelectedItem().(*listItem).internal.(*source.Chapter)

			if chapter.Volume == "" {
				break
			}

			for _, item := range b.chaptersC.Items() {
				item := item.(*listItem)
				if item.internal.(*source.Chapter).Volume == chapter.Volume {
					item.toggleMark()
					if item.marked {
						b.selectedChapters[item.internal.(*source.Chapter)] = struct{}{}
					} else {
						delete(b.selectedChapters, item.internal.(*source.Chapter))
					}
				}
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

		if b.chaptersToDownload.Len() == 0 {
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
			err := open.Start(lo.Must(b.currentDownloadingChapter.Manga.Path(false)))
			if err != nil {
				b.raiseError(err)
			}
		case key.Matches(msg, b.keymap.redownloadFailed):
			if len(b.failedChapters) == 0 {
				break
			}

			b.chaptersToDownload = util.Stack[*source.Chapter]{}
			for _, chapter := range b.failedChapters {
				b.chaptersToDownload.Push(chapter)
			}
			b.failedChapters = make([]*source.Chapter, 0)
			b.succededChapters = make([]*source.Chapter, 0)
			b.newState(downloadState)
			return b, tea.Batch(b.startLoading(), b.downloadChapter(b.chaptersToDownload.Pop()), b.waitForChapterDownload(), b.progressC.SetPercent(0))
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
