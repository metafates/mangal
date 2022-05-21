package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/metafates/mangai/api/downloader"
	"github.com/metafates/mangai/api/scraper"
	"log"
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

type progressInfoMsg progressInfo
type subProgressInfoMsg downloader.ChapterDownloadInfo
type sourceResponseMsg []*scraper.URL

func searchMangaAsync(b Bubble) tea.Cmd {
	return func() tea.Msg {
		goroutinesDone := 0
		var found []*scraper.URL
		for _, src := range b.config.UsedSources {
			src := src
			go func() {
				mangas, err := src.Mangas(b.input.Value())
				goroutinesDone++

				defer func() {
					if goroutinesDone == len(b.config.UsedSources) {
						b.sub <- found
					}
				}()

				if err != nil {
					return
				}

				found = append(found, mangas...)
			}()
		}

		return nil
	}
}

func getChaptersAsync(sub chan []*scraper.URL, url scraper.URL) tea.Cmd {
	return func() tea.Msg {
		chapters, err := url.Source.Chapters(url)
		if err != nil {
			sub <- []*scraper.URL{}
		}
		sub <- chapters
		return nil
	}
}

func startDownloaderAsync(b Bubble) tea.Cmd {
	return func() tea.Msg {
		items := b.chapters.Items()
		count := 0
		var (
			prevInfo string
			percent  float64
		)

		for index := range b.selected {
			chapter := items[index].(listItem)
			prevInfo = chapter.url.Info

			b.tick <- progressInfo{
				percent: percent,
				text:    prevInfo,
			}

			// TODO: Add error handling, maybe?
			_, _ = downloader.DownloadChapter(b.prevManga, chapter.url, b.subTick)
			count += 1
			percent = float64(count) / float64(len(b.selected))
		}

		b.tick <- progressInfo{
			percent: percent,
			text:    prevInfo,
		}

		return nil
	}
}

func waitForSourceResponse(sub chan []*scraper.URL) tea.Cmd {
	return func() tea.Msg {
		return sourceResponseMsg(<-sub)
	}
}

func waitForDownloaderResponse(sub chan progressInfo) tea.Cmd {
	return func() tea.Msg {
		return progressInfoMsg(<-sub)
	}
}

func waitForDownloaderSubResponse(sub chan downloader.ChapterDownloadInfo) tea.Cmd {
	return func() tea.Msg {
		return subProgressInfoMsg(<-sub)
	}
}

func (b Bubble) handleSearchState(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	k := b.keys[searchState]
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.resize(msg.Width, msg.Height)
		return b, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, k.Quit), key.Matches(msg, k.Back):
			return b, tea.Quit
		case key.Matches(msg, k.Confirm):
			b.state = spinnerState

			cmds = append(
				cmds,
				tea.Batch(b.spinner.Tick, searchMangaAsync(b), waitForSourceResponse(b.sub)),
			)
		}
	}

	b.input, cmd = b.input.Update(msg)
	cmds = append(cmds, cmd)
	return b, tea.Batch(cmds...)
}

func (b Bubble) handleSpinnerState(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	k := b.keys[spinnerState]
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.resize(msg.Width, msg.Height)
		return b, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, k.Quit):
			return b, tea.Quit
		case key.Matches(msg, k.Back):
			b.state = searchState
			// Ignore previous resolvers
			b.sub = make(chan []*scraper.URL)
		}
	case sourceResponseMsg:
		b.state = mangaSelectState
		var items []list.Item
		for _, url := range msg {
			items = append(items, listItem{url: *url})
		}
		cmds = append(cmds, b.manga.SetItems(items))
	}

	b.spinner, cmd = b.spinner.Update(msg)
	cmds = append(cmds, cmd)
	return b, tea.Batch(cmds...)
}

func (b Bubble) handleMangaSelectState(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	k := b.keys[mangaSelectState]
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.resize(msg.Width, msg.Height)
		return b, nil
	case sourceResponseMsg:
		b.state = chaptersSelectState
		var items []list.Item
		for _, url := range msg {
			items = append(items, listItem{url: *url})
		}
		b.manga.StopSpinner()

		cmds = append(cmds, b.chapters.SetItems(items), b.manga.NewStatusMessage(""))
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, k.Quit):
			return b, tea.Quit
		case key.Matches(msg, k.Back):
			b.state = searchState
			b.manga.Select(0)
			return b, nil
		case key.Matches(msg, k.Select), key.Matches(msg, k.Confirm):
			selected := b.manga.SelectedItem()

			if selected == nil {
				return b, nil
			}

			item, ok := selected.(listItem)

			b.prevManga = item.Title()

			if !ok {
				log.Fatal("Unknown manga is selected")
			}

			cmds = append(
				cmds,
				tea.Batch(
					getChaptersAsync(b.sub, item.url),
					waitForSourceResponse(b.sub),
					b.manga.StartSpinner(),
					b.manga.NewStatusMessage("Loading...")),
			)

		}
	}

	b.manga, cmd = b.manga.Update(msg)
	cmds = append(cmds, cmd)
	return b, tea.Batch(cmds...)
}

func (b Bubble) handleChaptersSelectState(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	k := b.keys[chaptersSelectState]
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.resize(msg.Width, msg.Height)
		return b, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, k.Quit):
			return b, tea.Quit

		case key.Matches(msg, k.Back):
			b.state = mangaSelectState
			b.selected = map[int]interface{}{}
			return b, nil

		case key.Matches(msg, k.Confirm):
			b.state = promptState
			return b, nil

		case key.Matches(msg, k.SelectAll):
			items := b.chapters.Items()
			for i, item := range items {
				l := item.(listItem)
				l.mark()

				// Toggle selected
				if _, exists := b.selected[i]; exists {
					delete(b.selected, i)
				} else {
					b.selected[i] = nil
				}

				items[i] = l
			}

			cmd = b.chapters.SetItems(items)
			return b, cmd

		case key.Matches(msg, k.Select):
			item, ok := b.chapters.SelectedItem().(listItem)
			index := b.chapters.Index()

			if !ok {
				log.Fatal("Unknown chapter is selected")
			}

			// Toggle selected
			if _, exists := b.selected[index]; exists {
				delete(b.selected, index)
			} else {
				b.selected[index] = nil
			}

			item.mark()
			cmd = b.chapters.SetItem(index, item)
			return b, cmd
		}
	}

	b.chapters, cmd = b.chapters.Update(msg)
	cmds = append(cmds, cmd)
	return b, tea.Batch(cmds...)
}

func (b Bubble) handlePromptState(msg tea.Msg) (tea.Model, tea.Cmd) {
	k := b.keys[promptState]

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.resize(msg.Width, msg.Height)
		return b, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, k.Back):
			b.state = chaptersSelectState
			return b, nil
		case key.Matches(msg, k.Quit):
			return b, tea.Quit
		case key.Matches(msg, k.Confirm):
			b.state = progressState
			return b, tea.Batch(startDownloaderAsync(b), waitForDownloaderResponse(b.tick), waitForDownloaderSubResponse(b.subTick))
		}
	}
	return b, nil
}

func (b Bubble) handleProgressState(msg tea.Msg) (tea.Model, tea.Cmd) {
	k := b.keys[progressState]

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.resize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, k.Quit):
			return b, tea.Quit
		}

	case progressInfoMsg:
		info := progressInfo(msg)

		if info.percent == 1 {
			b.state = exitPromptState
			return b, nil
		}

		cmd := b.progress.SetPercent(info.percent)
		b.prevChapter = info.text
		return b, tea.Batch(cmd, waitForDownloaderResponse(b.tick), waitForDownloaderSubResponse(b.subTick))

	case subProgressInfoMsg:
		info := downloader.ChapterDownloadInfo(msg)
		cmd := b.subProgress.SetPercent(info.Percent)
		b.prevPanel = info.Panel
		b.converting = info.Converting
		return b, tea.Batch(cmd, waitForDownloaderResponse(b.tick), waitForDownloaderSubResponse(b.subTick))

	case progress.FrameMsg:
		progressModel, cmd1 := b.progress.Update(msg)
		subProgressModel, cmd2 := b.subProgress.Update(msg)
		b.progress = progressModel.(progress.Model)
		b.subProgress = subProgressModel.(progress.Model)
		return b, tea.Batch(cmd1, cmd2)
	}

	return b, nil
}

func (b Bubble) handleExitPromptState(msg tea.Msg) (tea.Model, tea.Cmd) {
	k := b.keys[exitPromptState]
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.resize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, k.Back):
			b.state = chaptersSelectState
			return b, nil
		case key.Matches(msg, k.Quit):
			return b, tea.Quit
		}
	}

	return b, nil
}
