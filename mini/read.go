package mini

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/converter"
	"github.com/metafates/mangal/history"
	"github.com/metafates/mangal/log"
	"github.com/metafates/mangal/source"
	"github.com/samber/lo"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/viper"
	"os"
	"time"
)

func read() error {
	src, err := selectSource()
	if err != nil {
		return err
	}

	mangas, err := searchMangas(src)
	if err != nil {
		return err
	}

	manga, err := selectManga(mangas)
	if err != nil {
		return err
	}

	chapters, err := src.ChaptersOf(manga)
	if err != nil {
		return err
	}

	selected, err := selectChapter(chapters, 0)
	if err != nil {
		return err
	}

	return readChapter(src, selected)
}

func readChapter(src source.Source, chapter *source.Chapter) error {
	defer func() {
		if viper.GetBool(config.HistorySaveOnRead) {
			go func() {
				err := history.Save(chapter)
				if err != nil {
					log.Error(err)
				}
			}()
		}
	}()

	if viper.GetBool(config.ReaderReadInBrowser) {
		return open.Start(chapter.URL)
	}

	pages, err := src.PagesOf(chapter)
	if err != nil {
		return err
	}

	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.Suffix = fmt.Sprintf(" Downloading %d pages", len(pages))
	s.FinalMSG = finalMSG()
	lo.Must0(s.Color("bold", "magenta"))
	s.Start()
	err = chapter.DownloadPages()
	if err != nil {
		return err
	}

	s.Suffix = " Converting pages"

	conv, err := converter.Get(viper.GetString(config.FormatsUse))
	if err != nil {
		return err
	}

	path, err := conv.SaveTemp(chapter)
	if err != nil {
		return err
	}

	if reader := viper.GetString(config.ReaderName); reader != "" {
		s.Suffix = fmt.Sprintf(" Opening \"%s\"", reader)
		err = open.StartWith(path, reader)
	} else {
		s.Suffix = " Opening"
		err = open.Start(path)
	}

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	s.Stop()

	return nil
}
