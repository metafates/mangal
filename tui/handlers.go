package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/converter"
	"github.com/metafates/mangal/history"
	"github.com/metafates/mangal/source"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/viper"
)

func (b *statefulBubble) searchManga(query string) tea.Cmd {
	return func() tea.Msg {
		mangas, err := b.selectedSource.Search(query)
		if err != nil {
			b.errorChannel <- err
		} else {
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
		chapters, err := b.selectedSource.ChaptersOf(manga)
		if err != nil {
			b.errorChannel <- err
		} else {
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
		b.progressStatus = "Gettings pages"
		b.currentDownloadingChapter = chapter
		pages, err := b.selectedSource.PagesOf(chapter)
		if err != nil {
			b.errorChannel <- err
			return nil
		}

		b.progressStatus = fmt.Sprintf("Downloading %d pages", len(pages))
		err = chapter.DownloadPages()
		if err != nil {
			b.errorChannel <- err
			return nil
		}

		conv, err := converter.Get(viper.GetString(config.FormatsUse))
		if err != nil {
			b.errorChannel <- err
			return nil
		}

		b.progressStatus = fmt.Sprintf("Converting %d pages to %s", len(pages), viper.GetString(config.FormatsUse))
		path, err := conv.SaveTemp(chapter)
		if err != nil {
			b.errorChannel <- err
			return nil
		}

		if reader := viper.GetString(config.ReaderName); reader != "" {
			b.progressStatus = fmt.Sprintf("Opening %s", reader)
			err = open.RunWith(reader, path)
			if err != nil {
				b.errorChannel <- err
				return nil
			}
		} else {
			b.progressStatus = "Opening"
			err = open.Run(path)
			if err != nil {
				b.errorChannel <- err
				return nil
			}
		}

		_ = history.Save(chapter)
		b.progressStatus = "Done"
		b.chapterReadChannel <- struct{}{}

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
		b.progressStatus = "Getting pages"
		pages, err := b.selectedSource.PagesOf(chapter)
		if err != nil {
			b.errorChannel <- err
			return nil
		}

		b.progressStatus = fmt.Sprintf("Downloading %d pages", len(pages))
		err = chapter.DownloadPages()
		if err != nil {
			b.errorChannel <- err
			return nil
		}

		b.progressStatus = fmt.Sprintf("Converting %d pages to %s", len(pages), viper.GetString(config.FormatsUse))
		conv, err := converter.Get(viper.GetString(config.FormatsUse))
		if err != nil {
			b.errorChannel <- err
			return nil
		}

		path, err := conv.Save(chapter)
		if err != nil {
			b.errorChannel <- err
			return nil
		}

		b.progressStatus = "Downloaded"
		b.chapterDownloadChannel <- struct{}{}
		b.lastDownloadedChapterPath = path

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
