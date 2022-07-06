package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/metafates/mangal/cleaner"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/downloader"
	"github.com/metafates/mangal/history"
	"github.com/metafates/mangal/scraper"
	"github.com/metafates/mangal/util"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/exp/slices"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type mangaSearchDoneMsg []*scraper.URL

// initMangaSearch initializes the manga search
func (b *Bubble) initMangaSearch(query string) tea.Cmd {
	return func() tea.Msg {
		var (
			manga []*scraper.URL
			wg    sync.WaitGroup
		)

		wg.Add(len(config.UserConfig.Scrapers))

		for _, s := range config.UserConfig.Scrapers {
			go func(s *scraper.Scraper) {
				defer wg.Done()

				m, err := s.SearchManga(query)

				if err == nil {
					manga = append(manga, m...)
				}
			}(s)
		}

		wg.Wait()
		b.mangaChan <- manga

		return nil
	}
}

// waitForMangaSearchCompletion waits for the manga search to finish
func (b *Bubble) waitForMangaSearchCompletion() tea.Cmd {
	return func() tea.Msg {
		return mangaSearchDoneMsg(<-b.mangaChan)
	}
}

type chaptersGetDoneMsg []*scraper.URL
type chapterDownloadProgressMsg downloader.ChapterDownloadProgress

// initChaptersGet initializes the chapters get
func (b *Bubble) initChaptersGet(manga *scraper.URL) tea.Cmd {
	return func() tea.Msg {
		chapters, err := manga.Scraper.GetChapters(manga)

		if err != nil {
			// set to empty list if error occured and notify user
			b.chaptersChan <- make([]*scraper.URL, 0)
			return b.chaptersList.NewStatusMessage("Error occured while fetching chapters")()
		}

		if config.UserConfig.Anilist.Enabled {
			// cache result
			config.UserConfig.Anilist.Client.ToAnilistURL(manga)
		}

		b.chaptersChan <- chapters
		return nil
	}
}

// waitForChaptersGetCompletion waits for the chapters get to finish
func (b *Bubble) waitForChaptersGetCompletion() tea.Cmd {
	return func() tea.Msg {
		return chaptersGetDoneMsg(<-b.chaptersChan)
	}
}

// waitForChapterDownloadProgress waits for the chapter download progress to finish
func (b *Bubble) waitForChapterDownloadProgress() tea.Cmd {
	return func() tea.Msg {
		return chapterDownloadProgressMsg(<-b.chapterPagesProgressChan)
	}
}

type chaptersDownloadProgressMsg downloader.ChaptersDownloadProgress

// initChaptersDownload initializes the chapters download
func (b *Bubble) initChaptersDownload(chapters []*scraper.URL) tea.Cmd {
	return func() tea.Msg {
		var (
			failed    []*scraper.URL
			succeeded []string
			total     = len(chapters)
			path      string
			err       error
		)

		sort.Slice(chapters, func(i int, j int) bool {
			return chapters[i].Index < chapters[j].Index
		})

		// Download chapters
		for i, chapter := range chapters {
			b.chaptersProgressChan <- downloader.ChaptersDownloadProgress{
				Failed:    failed,
				Succeeded: succeeded,
				Total:     total,
				Proceeded: util.Max(i-1, 0),
				Current:   chapter,
				Done:      false,
			}

			path, err = downloader.DownloadChapter(chapter, b.chapterPagesProgressChan, false)
			if err == nil {
				if config.UserConfig.Anilist.MarkDownloaded {
					_ = config.UserConfig.Anilist.Client.MarkChapter(chapter.Relation, chapter.Index)
				}

				// use path instead of the chapter name since it is used to get manga folder later
				succeeded = append(succeeded, path)
			} else {
				failed = append(failed, chapter)
			}
		}

		// If epub file was used, create it
		if downloader.EpubFile != nil {
			downloader.EpubFile.SetAuthor(chapters[0].Scraper.Source.Base)
			if err := downloader.EpubFile.Write(path); err != nil {
				log.Fatal("Error while making epub. Please, try again")
			}

			// Close epub file
			downloader.EpubFile = nil
		}

		b.chaptersProgressChan <- downloader.ChaptersDownloadProgress{
			Failed:    failed,
			Succeeded: succeeded,
			Total:     total,
			Proceeded: total,
			Current:   nil,
			Done:      true,
		}

		return nil
	}
}

// waitForChaptersDownloadProgress waits for the chapters download progress to finish
func (b *Bubble) waitForChaptersDownloadProgress() tea.Cmd {
	return func() tea.Msg {
		return chaptersDownloadProgressMsg(<-b.chaptersProgressChan)
	}
}

type chapterDownloadedToReadMsg downloader.ChaptersDownloadProgress

// initChapterDownloadedToRead initializes the chapter downloaded to read
func (b *Bubble) initChapterDownloadToRead(chapter *scraper.URL) tea.Cmd {
	return func() tea.Msg {
		var (
			failed    []*scraper.URL
			succeeded []string
		)

		b.chaptersProgressChan <- downloader.ChaptersDownloadProgress{
			Current:   chapter,
			Done:      false,
			Failed:    failed,
			Succeeded: succeeded,
			Total:     1,
			Proceeded: 0,
		}

		path, err := downloader.DownloadChapter(chapter, b.chapterPagesProgressChan, true)

		if err != nil {
			failed = append(failed, chapter)
		} else {
			// Mark chapter as read
			if config.UserConfig.Anilist.Enabled {
				go func() {
					_ = config.UserConfig.Anilist.Client.MarkChapter(chapter.Relation, chapter.Index)
				}()
			}

			succeeded = append(succeeded, path)
		}

		b.chaptersProgressChan <- downloader.ChaptersDownloadProgress{
			Current:   nil,
			Done:      true,
			Failed:    failed,
			Succeeded: succeeded,
			Total:     1,
			Proceeded: 1,
		}

		return nil
	}
}

// waitForChapterToReadDownloaded waits for the chapter to read download to finish
func (b *Bubble) waitForChapterToReadDownloaded() tea.Cmd {
	return func() tea.Msg {
		return chapterDownloadedToReadMsg(<-b.chaptersProgressChan)
	}
}

func (b *Bubble) switchToChapters(msg chaptersGetDoneMsg, chapterIndex int) {
	b.setState(ChaptersState)
	var anilistManga *scraper.AnilistURL

	if len(msg) > 0 {
		manga := msg[0].Relation
		if config.UserConfig.Anilist.Enabled {
			anilistManga = config.UserConfig.Anilist.Client.ToAnilistURL(manga)
		}
		b.chaptersList.Title = "Chapters - " + util.PrettyTrim(manga.Info, 30)
	} else {
		b.chaptersList.Title = "Chapters"
	}

	var items []list.Item

	// Sort according to chapter index, in ascending order
	sort.Slice(msg, func(i, j int) bool {
		return msg[i].Index < msg[j].Index
	})

	for _, url := range msg {
		items = append(items, &listItem{url: url})
	}

	var cmds []tea.Cmd

	cmds = append(cmds, b.chaptersList.SetItems(items))
	if anilistManga != nil {
		cmds = append(cmds, b.chaptersList.NewStatusMessage("AL: "+util.PrettyTrim(anilistManga.Title, 25)))
	}
	b.mangaList.StopSpinner()
	b.ResumeList.StopSpinner()
	b.chaptersList.Select(chapterIndex)
}

type selectIndexMsg int

func (b *Bubble) selectIndex(index int) tea.Cmd {
	return func() tea.Msg {
		return selectIndexMsg(index)
	}
}

type resumeMsg *history.Entry

func (b *Bubble) resumeFromHistory(entry *history.Entry) tea.Cmd {
	return func() tea.Msg {
		return resumeMsg(entry)
	}
}

func (b *Bubble) handleResumeState(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keyMap.Back):
			if slices.Contains([]list.FilterState{list.FilterApplied, list.Filtering}, b.ResumeList.FilterState()) {
				b.ResumeList.ResetFilter()
				return b, nil
			}

			return b, tea.Quit
		case key.Matches(msg, b.keyMap.Quit):
			if b.ResumeList.FilterState() == list.Filtering {
				break
			}

			return b, tea.Quit
		case key.Matches(msg, b.keyMap.Open):
			if b.ResumeList.FilterState() == list.Filtering {
				break
			}

			selected, ok := b.ResumeList.SelectedItem().(*history.Entry)

			if ok {
				_ = open.Start(selected.Manga.Address)
			}
		case key.Matches(msg, b.keyMap.Select), key.Matches(msg, b.keyMap.Confirm):
			if key.Matches(msg, b.keyMap.Select) && b.ResumeList.FilterState() == list.Filtering {
				break
			}

			// do not select twice while loading
			if b.loading {
				return b, nil
			}

			selected, ok := b.ResumeList.SelectedItem().(*history.Entry)

			if ok {
				b.loading = true
				HistoryChapterIndex = selected.Chapter.Index - 1
				return b, tea.Batch(
					b.ResumeList.StartSpinner(),
					b.initChaptersGet(selected.Manga),
					b.waitForChaptersGetCompletion(),
					b.ResumeList.StartSpinner(),
				)
			}
		}
	case chaptersGetDoneMsg:
		b.loading = false
		b.ResumeList.ResetFilter()
		b.switchToChapters(msg, HistoryChapterIndex)
		return b, nil
	}

	*b.ResumeList, cmd = b.ResumeList.Update(msg)

	return b, cmd
}

// handleSearchState handles the search state
func (b *Bubble) handleSearchState(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keyMap.Back):
			return b, tea.Quit
		case key.Matches(msg, b.keyMap.Confirm) && b.input.Value() != "":
			b.setState(LoadingState)
			return b, tea.Batch(
				b.spinner.Tick,
				b.initMangaSearch(b.input.Value()),
				b.waitForMangaSearchCompletion(),
			)
		}
	}

	*b.input, cmd = b.input.Update(msg)
	return b, cmd
}

// handleLoadingState handles the loading state
func (b *Bubble) handleLoadingState(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case mangaSearchDoneMsg:
		b.setState(MangaState)
		b.mangaList.Title = "Manga - " + util.PrettyTrim(strings.TrimSpace(b.input.Value()), 30)

		var items = make([]list.Item, len(msg))

		for i, url := range msg {
			items[i] = &listItem{url: url}
		}

		cmd = b.mangaList.SetItems(items)
		return b, cmd
	}

	*b.spinner, cmd = b.spinner.Update(msg)
	return b, cmd
}

// handleMangaState handles the manga state
func (b *Bubble) handleMangaState(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keyMap.Back):
			b.loading = false
			b.mangaList.StopSpinner()
			b.mangaList.Select(0)
			b.setState(SearchState)
			return b, b.chaptersList.NewStatusMessage("")

		case key.Matches(msg, b.keyMap.Quit):
			return b, tea.Quit
		case key.Matches(msg, b.keyMap.Open):
			item, ok := b.mangaList.SelectedItem().(*listItem)
			if ok {
				_ = open.Start(item.url.Address)
			}
		case b.loading:
			// Do nothing if the chapters are loading
			return b, nil
		case key.Matches(msg, b.keyMap.Confirm), key.Matches(msg, b.keyMap.Select):
			selected, ok := b.mangaList.SelectedItem().(*listItem)
			if ok {
				cmd = b.mangaList.StartSpinner()
				b.loading = true

				return b, tea.Batch(cmd, b.initChaptersGet(selected.url), b.waitForChaptersGetCompletion())
			}
		}
	case chaptersGetDoneMsg:
		b.loading = false
		b.switchToChapters(msg, 0)
		return b, cmd
	}

	*b.mangaList, cmd = b.mangaList.Update(msg)
	return b, cmd
}

// handleChaptersState handles the chapters state
func (b *Bubble) handleChaptersState(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keyMap.Quit):
			return b, tea.Quit
		case key.Matches(msg, b.keyMap.Back):
			// reset selected items
			b.chaptersList.Select(0)
			b.selectedChapters = make(map[int]interface{})

			if config.UserConfig.HistoryMode {
				b.setState(ResumeState)

				// update history entries
				h, err := history.ReadHistory()

				if err != nil {
					log.Fatal(err)
				}

				var items []list.Item

				for _, item := range h {
					items = append(items, item)
				}

				sort.Slice(items, func(i, j int) bool {
					return items[i].(*history.Entry).Manga.Info < items[j].(*history.Entry).Manga.Info
				})

				b.ResumeList.SetItems(items)
			} else {
				b.setState(MangaState)
			}

			b.chaptersList.NewStatusMessage("")

			return b, cmd
		case key.Matches(msg, b.keyMap.Open):
			item, ok := b.chaptersList.SelectedItem().(*listItem)
			if ok {
				_ = open.Start(item.url.Address)
			}
		case key.Matches(msg, b.keyMap.Read):
			chapter, ok := b.chaptersList.SelectedItem().(*listItem)

			if ok {
				b.setState(DownloadingState)

				if !config.UserConfig.IncognitoMode {
					go func() {
						_ = history.WriteHistory(chapter.url)
					}()
				}

				return b, tea.Batch(
					b.progress.SetPercent(0),
					b.spinner.Tick,
					b.initChapterDownloadToRead(chapter.url),
					b.waitForChapterToReadDownloaded(),
					b.waitForChapterDownloadProgress(),
				)
			}
		case key.Matches(msg, b.keyMap.Confirm) && len(b.selectedChapters) > 0:
			b.setState(ConfirmState)
			return b, nil
		case key.Matches(msg, b.keyMap.Select):
			item, ok := b.chaptersList.SelectedItem().(*listItem)
			if ok {
				index := b.chaptersList.Index()
				item.Select()

				if item.selected {
					b.selectedChapters[index] = nil
				} else {
					delete(b.selectedChapters, index)
				}

				return b, nil
			}
		case key.Matches(msg, b.keyMap.SelectAll):
			items := b.chaptersList.Items()

			for i, item := range items {
				it := item.(*listItem)
				it.Select()

				if it.selected {
					b.selectedChapters[i] = nil
				} else {
					delete(b.selectedChapters, i)
				}
			}
			return b, nil
		}
	}

	*b.chaptersList, cmd = b.chaptersList.Update(msg)
	return b, cmd
}

// handleConfirmPromptState handles the confirmation prompt state
func (b *Bubble) handleConfirmPromptState(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keyMap.Quit):
			return b, tea.Quit
		case key.Matches(msg, b.keyMap.Back):
			b.setState(ChaptersState)
			return b, nil
		case key.Matches(msg, b.keyMap.Confirm):
			b.setState(DownloadingState)

			var (
				chapters = make([]*scraper.URL, len(b.selectedChapters))
				iterIdx  int
			)

			items := b.chaptersList.Items()

			for index := range b.selectedChapters {
				chapters[iterIdx] = items[index].(*listItem).url
				iterIdx++
			}

			return b, tea.Batch(
				b.progress.SetPercent(0),
				b.spinner.Tick,
				b.initChaptersDownload(chapters),
				b.waitForChaptersDownloadProgress(),
				b.waitForChapterDownloadProgress(),
			)
		}
	}

	return b, nil
}

// handleDownloadingState handles the downloading state
func (b *Bubble) handleDownloadingState(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case chapterDownloadedToReadMsg:
		b.chaptersDownloadProgressInfo = downloader.ChaptersDownloadProgress(msg)

		if msg.Done {
			b.setState(ExitPromptState)

			if len(msg.Succeeded) != 0 {
				toRead := msg.Succeeded[0]

				if config.UserConfig.UseCustomReader {
					_ = open.StartWith(toRead, config.UserConfig.CustomReader)
				} else {
					_ = open.Start(toRead)
				}
			}

			return b, nil
		}

		cmd = b.progress.SetPercent(float64(len(msg.Succeeded)) / float64(msg.Total))

		return b, tea.Batch(cmd, b.waitForChapterToReadDownloaded(), b.waitForChapterDownloadProgress())
	case chapterDownloadProgressMsg:
		*b.spinner, cmd = b.spinner.Update(msg)
		b.chapterDownloadProgressInfo = downloader.ChapterDownloadProgress(msg)
		return b, tea.Batch(cmd, b.waitForChapterDownloadProgress(), b.waitForChaptersGetCompletion())
	case chaptersDownloadProgressMsg:
		b.chaptersDownloadProgressInfo = downloader.ChaptersDownloadProgress(msg)

		if msg.Done {
			b.setState(ExitPromptState)
			return b, nil
		}

		cmd = b.progress.SetPercent(float64(len(msg.Succeeded)) / float64(msg.Total))
		return b, tea.Batch(cmd, b.waitForChaptersDownloadProgress(), b.waitForChapterDownloadProgress())
	case progress.FrameMsg:
		var p tea.Model
		// ???? why progress.Update() returns tea.Model and not progress.Model?
		p, cmd = b.progress.Update(msg)
		*b.progress = p.(progress.Model)
		return b, cmd
	}

	*b.spinner, cmd = b.spinner.Update(msg)
	return b, cmd
}

// handleExitPromptState handles the exit prompt state
func (b *Bubble) handleExitPromptState(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keyMap.Quit):
			cleaner.RemoveTemp()
			return b, tea.Quit
		case key.Matches(msg, b.keyMap.Back):
			b.setState(ChaptersState)
			return b, nil
		case key.Matches(msg, b.keyMap.Retry):
			failed := b.chaptersDownloadProgressInfo.Failed

			if len(failed) > 0 {
				b.setState(DownloadingState)
				return b, tea.Batch(
					b.progress.SetPercent(0),
					b.spinner.Tick,
					b.initChaptersDownload(failed),
					b.waitForChaptersDownloadProgress(),
					b.waitForChapterDownloadProgress(),
				)
			}
		case key.Matches(msg, b.keyMap.Open):
			if paths := b.chaptersDownloadProgressInfo.Succeeded; len(paths) > 0 {
				_ = open.Start(filepath.Dir(paths[0]))
			}
		}
	}

	return b, nil
}

// Update handles the Bubble update
func (b *Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.resize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keyMap.ForceQuit):
			cleaner.RemoveTemp()
			return b, tea.Quit
		}
	}

	switch b.state {
	case ResumeState:
		return b.handleResumeState(msg)
	case SearchState:
		return b.handleSearchState(msg)
	case LoadingState:
		return b.handleLoadingState(msg)
	case MangaState:
		return b.handleMangaState(msg)
	case ChaptersState:
		return b.handleChaptersState(msg)
	case ConfirmState:
		return b.handleConfirmPromptState(msg)
	case DownloadingState:
		return b.handleDownloadingState(msg)
	case ExitPromptState:
		return b.handleExitPromptState(msg)
	}

	log.Fatal("Unknown state encountered")

	// Unreachable
	return b, nil
}
