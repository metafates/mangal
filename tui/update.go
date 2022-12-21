package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/metafates/mangal/anilist"
	"github.com/metafates/mangal/color"
	"github.com/metafates/mangal/history"
	"github.com/metafates/mangal/installer"
	key2 "github.com/metafates/mangal/key"
	"github.com/metafates/mangal/open"
	"github.com/metafates/mangal/provider"
	"github.com/metafates/mangal/query"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/style"
	"github.com/metafates/mangal/util"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/spf13/viper"
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
			onListBack := func(l *list.Model) tea.Cmd {
				l.ResetSelected()
				l.ResetFilter()

				return tea.Batch(cmd, l.NewStatusMessage(""))
			}

			switch b.state {
			case searchState:
				b.inputC.SetValue("")
			case chaptersState:
				if b.chaptersC.FilterState() != list.Unfiltered {
					b.chaptersC, cmd = b.chaptersC.Update(msg)
					return b, cmd
				}

				b.selectedChapters = make(map[*source.Chapter]struct{})
				cmd = onListBack(&b.chaptersC)
			case anilistSelectState:
				if b.anilistC.FilterState() != list.Unfiltered {
					b.anilistC, cmd = b.anilistC.Update(msg)
					return b, cmd
				}

				cmd = onListBack(&b.anilistC)
			case mangasState:
				if b.mangasC.FilterState() != list.Unfiltered {
					b.mangasC, cmd = b.mangasC.Update(msg)
					return b, cmd
				}

				cmd = onListBack(&b.mangasC)
			case historyState:
				if b.historyC.FilterState() != list.Unfiltered {
					b.historyC, cmd = b.historyC.Update(msg)
					return b, cmd
				}

				cmd = onListBack(&b.historyC)
			case sourcesState:
				if b.sourcesC.FilterState() != list.Unfiltered {
					b.sourcesC, cmd = b.sourcesC.Update(msg)
					return b, cmd
				}

				cmd = onListBack(&b.sourcesC)
			case scrapersInstallState:
				if b.scrapersInstallC.FilterState() != list.Unfiltered {
					b.scrapersInstallC, cmd = b.scrapersInstallC.Update(msg)
					return b, cmd
				}

				cmd = onListBack(&b.scrapersInstallC)
			}

			b.previousState()
			b.stopLoading()
			b.failedChapters = make([]*source.Chapter, 0)
			b.succededChapters = make([]*source.Chapter, 0)
			return b, cmd
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
	case anilistSelectState:
		return b.updateAnilistSelect(msg)
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
	case []*anilist.Manga:
		manga, err := anilist.FindClosest(b.selectedManga.Name)
		id := -1
		if err == nil {
			id = manga.ID
		}

		items := make([]list.Item, len(msg))
		var marked int
		for i, manga := range msg {
			if manga.ID == id {
				marked = i
			}

			items[i] = &listItem{
				internal: manga,
				marked:   manga.ID == id,
			}
		}

		cmd = b.anilistC.SetItems(items)
		b.newState(anilistSelectState)
		b.anilistC.Select(marked)
		return b, tea.Batch(cmd, b.stopLoading())
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
			items[i] = &listItem{internal: m}
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
			items[i] = &listItem{internal: c}
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
			defer b.newState(loadingState)

			if len(b.selectedProviders) == 0 {
				return b, tea.Batch(b.startLoading(), b.loadSources([]*provider.Provider{item.internal.(*provider.Provider)}), b.waitForSourcesLoaded())
			}

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
			go query.Remember(b.inputC.Value(), 1)
			return b, tea.Batch(b.searchManga(b.inputC.Value()), b.waitForMangas(), b.spinnerC.Tick)
		case key.Matches(msg, b.keymap.acceptSearchSuggestion) && b.searchSuggestion.IsPresent():
			b.inputC.SetValue(b.searchSuggestion.MustGet())
			b.searchSuggestion = mo.None[string]()
			b.inputC.SetCursor(len(b.inputC.Value()))
			return b, nil
		}
	}

	b.inputC, cmd = b.inputC.Update(msg)

	if b.inputC.Value() != "" {
		if suggestion, ok := query.Suggest(b.inputC.Value()).Get(); ok && suggestion != b.inputC.Value() {
			b.searchSuggestion = mo.Some(suggestion)
		} else {
			b.searchSuggestion = mo.None[string]()
		}
	} else if b.searchSuggestion.IsPresent() {
		b.searchSuggestion = mo.None[string]()
	}

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
			go query.Remember(m.Name, 2)
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

		if viper.GetBool(key2.TUIReverseChapters) {
			for i, c := range msg {
				items[len(msg)-i-1] = &listItem{internal: c}
			}
		} else {
			for i, c := range msg {
				items[i] = &listItem{internal: c}
			}
		}

		cmd = b.chaptersC.SetItems(items)
		b.newState(chaptersState)
		b.stopLoading()

		if viper.GetBool(key2.AnilistLinkOnMangaSelect) {
			return b, tea.Batch(cmd, b.fetchAndSetAnilist(b.selectedManga), b.waitForAnilistFetchAndSet())
		}

		return b, cmd
	}

	b.mangasC, cmd = b.mangasC.Update(msg)
	return b, cmd
}

func (b *statefulBubble) updateChapters(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case *anilist.Manga:
		cmd = b.chaptersC.NewStatusMessage(fmt.Sprintf(`Linked to %s %s`, style.Fg(color.Orange)(msg.Name()), style.Faint(msg.SiteURL)))
		return b, cmd
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
		case key.Matches(msg, b.keymap.anilistSelect):
			b.newState(loadingState)
			return b, tea.Batch(b.startLoading(), b.fetchAnilist(b.selectedManga), b.waitForAnilist())
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
					if !item.marked {
						b.selectedChapters[item.internal.(*source.Chapter)] = struct{}{}
					}
					item.marked = true
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
			} else if viper.GetBool(key2.TUIReadOnEnter) {
				if b.chaptersC.SelectedItem() == nil {
					break
				}

				chapter := b.chaptersC.SelectedItem().(*listItem).internal.(*source.Chapter)
				b.newState(readState)
				return b, tea.Batch(b.readChapter(chapter), b.waitForChapterRead(), b.startLoading())
			}
		}
	}

	b.chaptersC, cmd = b.chaptersC.Update(msg)
	return b, cmd
}

func (b *statefulBubble) updateAnilistSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case b.anilistC.FilterState() == list.Filtering:
			break
		case key.Matches(msg, b.keymap.openURL):
			if b.anilistC.SelectedItem() == nil {
				break
			}

			m, _ := b.anilistC.SelectedItem().(*listItem).internal.(*anilist.Manga)
			err := open.Start(m.SiteURL)
			if err != nil {
				b.raiseError(err)
			}
		case key.Matches(msg, b.keymap.confirm, b.keymap.selectOne):
			if b.anilistC.SelectedItem() == nil {
				break
			}

			manga := b.anilistC.SelectedItem().(*listItem).internal.(*anilist.Manga)
			err := anilist.SetRelation(b.selectedManga.Name, manga)
			if err != nil {
				b.raiseError(err)
				break
			}

			b.previousState()
			cmd = b.chaptersC.NewStatusMessage(fmt.Sprintf(`Linked to %s %s`, style.Fg(color.Orange)(manga.Name()), style.Faint(manga.SiteURL)))
			return b, cmd
		}
	}

	b.anilistC, cmd = b.anilistC.Update(msg)
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
			err := open.StartWith(
				lo.Must(b.currentDownloadingChapter.Manga.Path(false)),
				viper.GetString(key2.ReaderFolder),
			)

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
