package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/metafates/mangal/anilist"
	"github.com/metafates/mangal/color"
	"github.com/metafates/mangal/downloader"
	"github.com/metafates/mangal/installer"
	"github.com/metafates/mangal/key"
	"github.com/metafates/mangal/log"
	"github.com/metafates/mangal/provider"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/style"
	"github.com/metafates/mangal/util"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
	"strings"
	"sync"
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

		var items = make([]list.Item, len(scrapers))
		for i, s := range scrapers {
			items[i] = &listItem{
				internal: s,
			}
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

func (b *statefulBubble) loadSources(ps []*provider.Provider) tea.Cmd {
	return func() tea.Msg {
		var (
			sources = make([]source.Source, len(ps))
			wg      = sync.WaitGroup{}
			err     error
		)

		wg.Add(len(ps))
		for i, p := range ps {
			go func(i int, p *provider.Provider) {
				defer wg.Done()

				if err != nil {
					return
				}

				log.Info("loading source " + p.ID)
				b.progressStatus = "Initializing source"
				var s source.Source
				s, err = p.CreateSource()

				if err != nil {
					log.Error(err)
					b.errorChannel <- err
					return
				}

				log.Info("source " + p.ID + " loaded")
				sources[i] = s
			}(i, p)
		}

		wg.Wait()

		b.sourcesLoadedChannel <- sources

		return nil
	}
}

func (b *statefulBubble) waitForSourcesLoaded() tea.Cmd {
	return func() tea.Msg {
		select {
		case res := <-b.sourcesLoadedChannel:
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
		b.progressStatus = fmt.Sprintf("Searching among %s", util.Quantify(len(b.selectedSources), "source", "sources"))

		var mangas = make([]*source.Manga, 0)

		wg := sync.WaitGroup{}
		wg.Add(len(b.selectedSources))
		for _, s := range b.selectedSources {
			go func(s source.Source) {
				defer wg.Done()
				sourceMangas, err := s.Search(query)

				if err != nil {
					log.Error(err)
					b.errorChannel <- err
				}

				log.Infof("found %s from source %s", util.Quantify(len(sourceMangas), "manga", "mangas"), s.Name())
				mangas = append(mangas, sourceMangas...)
			}(s)
		}

		wg.Wait()

		log.Infof("found %d mangas from %d sources", len(mangas), len(b.selectedSources))

		b.foundMangasChannel <- mangas

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
		chapters, err := manga.Source.ChaptersOf(manga)
		if err != nil {
			log.Error(err)
			b.errorChannel <- err
		} else {
			log.Infof("found %s", util.Quantify(len(chapters), "chapter", "chapters"))
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
		err := downloader.Read(chapter, func(s string) {
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
		_, err := downloader.Download(chapter, func(s string) {
			b.progressStatus = s
		})

		if err != nil {
			if viper.GetBool(key.DownloaderStopOnError) {
				b.errorChannel <- err
			} else {
				b.failedChapters = append(b.failedChapters, chapter)
				b.chapterDownloadChannel <- struct{}{}
			}
		} else {
			b.succededChapters = append(b.succededChapters, chapter)
			b.chapterDownloadChannel <- struct{}{}
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

func (b *statefulBubble) fetchAndSetAnilist(manga *source.Manga) tea.Cmd {
	return func() tea.Msg {
		alManga, err := anilist.FindClosest(manga.Name)
		if err != nil {
			// this error is not that important, we can ignore t
			log.Warn(err)
		} else {
			b.closestAnilistMangaChannel <- alManga
		}

		return nil
	}
}

func (b *statefulBubble) waitForAnilistFetchAndSet() tea.Cmd {
	return func() tea.Msg {
		return <-b.closestAnilistMangaChannel
	}
}

func (b *statefulBubble) fetchAnilist(manga *source.Manga) tea.Cmd {
	return func() tea.Msg {
		log.Info("fetching anilist for " + manga.Name)
		b.progressStatus = fmt.Sprintf("Fetching anilist for %s", style.Fg(color.Purple)(manga.Name))
		mangas, err := anilist.SearchByName(manga.Name)
		if err != nil {
			log.Error(err)
			b.errorChannel <- err
		} else {
			log.Infof("found %s", util.Quantify(len(mangas), "manga", "mangas"))
			b.fetchedAnilistMangasChannel <- mangas
		}

		return nil
	}
}

func (b *statefulBubble) waitForAnilist() tea.Cmd {
	return func() tea.Msg {
		select {
		case found := <-b.fetchedAnilistMangasChannel:
			return found
		case err := <-b.errorChannel:
			b.lastError = err
			return err
		}
	}
}

func (b *statefulBubble) selectChapterBy(f func(chapter *source.Chapter) bool) tea.Cmd {
	return func() tea.Msg {
		for i, item := range b.chaptersC.Items() {
			chapter := item.(*listItem).internal.(*source.Chapter)
			if f(chapter) {
				b.chaptersC.Select(i)
				return nil
			}
		}

		return nil
	}
}
