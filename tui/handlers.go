package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/downloader"
	"github.com/metafates/mangal/installer"
	"github.com/metafates/mangal/log"
	"github.com/metafates/mangal/provider"
	"github.com/metafates/mangal/source"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
	"strings"
)

func (b *statefulBubble) loadScrapers() tea.Cmd {
	return func() tea.Msg {
		b.progressStatus = "Loading scrapers"
		scrapers, err := installer.Scrapers()
		if err != nil {
			log.Error(err)
			b.errorChannel <- err
			return nil
		}
		b.progressStatus = "Scrapers Loaded"

		slices.SortFunc(scrapers, func(a, b *installer.Scraper) bool {
			return strings.Compare(a.Name, b.Name) < 0
		})

		var items []list.Item
		for _, s := range scrapers {
			items = append(items, &listItem{
				title:       s.Name,
				description: s.GithubURL(),
				internal:    s,
			})
		}

		cmd := b.scrapersInstallC.SetItems(items)
		b.scrapersLoadedChannel <- scrapers
		return cmd
	}
}

func (b *statefulBubble) waitForScrapersLoaded() tea.Cmd {
	return func() tea.Msg {
		select {
		case res := <-b.scrapersLoadedChannel:
			return res
		case err := <-b.errorChannel:
			b.lastError = err
			return err
		}
	}
}

func (b *statefulBubble) installScraper(s *installer.Scraper) tea.Cmd {
	return func() tea.Msg {
		b.progressStatus = fmt.Sprintf("Installing %s", s.Name)
		err := s.Install()
		if err != nil {
			log.Error(err)
			b.errorChannel <- err
		} else {
			log.Info("scraper " + s.Name + " installed")
			b.scraperInstalledChannel <- s
		}

		return nil
	}
}

func (b *statefulBubble) waitForScraperInstallation() tea.Cmd {
	return func() tea.Msg {
		select {
		case res := <-b.scraperInstalledChannel:
			return res
		case err := <-b.errorChannel:
			b.lastError = err
			return err
		}
	}
}

func (b *statefulBubble) loadSource(p *provider.Provider) tea.Cmd {
	return func() tea.Msg {
		log.Info("loading source " + p.ID)
		b.progressStatus = "Initializing source"
		s, err := p.CreateSource()

		if err != nil {
			log.Error(err)
			b.errorChannel <- err
		} else {
			log.Info("source " + p.ID + " loaded")
			b.sourceLoadedChannel <- s
		}

		return nil
	}
}

func (b *statefulBubble) waitForSourceLoaded() tea.Cmd {
	return func() tea.Msg {
		select {
		case res := <-b.sourceLoadedChannel:
			return res
		case err := <-b.errorChannel:
			b.lastError = err
			return err
		}
	}
}

func (b *statefulBubble) searchManga(query string) tea.Cmd {
	return func() tea.Msg {
		log.Info("searching for " + query)
		b.progressStatus = "Searching"
		mangas, err := b.selectedSource.Search(query)
		if err != nil {
			log.Error(err)
			b.errorChannel <- err
		} else {
			log.Info("found " + fmt.Sprintf("%d", len(mangas)) + " mangas")
			b.foundMangasChannel <- mangas
		}

		return nil
	}
}

func (b *statefulBubble) waitForMangas() tea.Cmd {
	return func() tea.Msg {
		select {
		case found := <-b.foundMangasChannel:
			return found
		case err := <-b.errorChannel:
			b.lastError = err
			return err
		}
	}
}

func (b *statefulBubble) getChapters(manga *source.Manga) tea.Cmd {
	return func() tea.Msg {
		log.Info("getting chapters of " + manga.Name)
		chapters, err := b.selectedSource.ChaptersOf(manga)
		if err != nil {
			log.Error(err)
			b.errorChannel <- err
		} else {
			log.Info("found " + fmt.Sprintf("%d", len(chapters)) + " chapters")
			b.foundChaptersChannel <- chapters
		}

		return nil
	}
}

func (b *statefulBubble) waitForChapters() tea.Cmd {
	return func() tea.Msg {
		select {
		case found := <-b.foundChaptersChannel:
			return found
		case err := <-b.errorChannel:
			b.lastError = err
			return err
		}
	}
}

func (b *statefulBubble) readChapter(chapter *source.Chapter) tea.Cmd {
	return func() tea.Msg {
		b.currentDownloadingChapter = chapter
		err := downloader.Read(b.selectedSource, chapter, func(s string) {
			b.progressStatus = s
		})

		if err != nil {
			b.errorChannel <- err
		} else {
			b.chapterReadChannel <- struct{}{}
		}

		return nil
	}
}

func (b *statefulBubble) waitForChapterRead() tea.Cmd {
	return func() tea.Msg {
		select {
		case res := <-b.chapterReadChannel:
			return res
		case err := <-b.errorChannel:
			b.lastError = err
			return err
		}
	}
}

func (b *statefulBubble) downloadChapter(chapter *source.Chapter) tea.Cmd {
	return func() tea.Msg {
		b.currentDownloadingChapter = chapter
		path, err := downloader.Download(b.selectedSource, chapter, func(s string) {
			b.progressStatus = s
		})

		if err != nil {
			if viper.GetBool(config.DownloaderStopOnError) {
				b.errorChannel <- err
			} else {
				b.failedChapters = append(b.failedChapters, chapter)
				b.chapterDownloadChannel <- struct{}{}
			}
		} else {
			b.succededChapters = append(b.succededChapters, chapter)
			b.chapterDownloadChannel <- struct{}{}
			b.lastDownloadedChapterPath = path
		}

		return nil
	}
}

func (b *statefulBubble) waitForChapterDownload() tea.Cmd {
	return func() tea.Msg {
		select {
		case res := <-b.chapterDownloadChannel:
			return res
		case err := <-b.errorChannel:
			b.lastError = err
			return err
		}
	}
}
