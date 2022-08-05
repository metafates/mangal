package mini

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/converter"
	"github.com/metafates/mangal/history"
	"github.com/metafates/mangal/style"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"os"
	"time"
)

func download() error {
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

	selected, err := selectChapters(chapters)
	if err != nil {
		return err
	}

	var counter int
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	lo.Must0(s.Color("bold", "magenta"))
	s.Suffix = " Starting..."
	s.FinalMSG = finalMSG()
	s.Start()

	conv, err := converter.Get(viper.GetString(config.FormatsUse))
	if err != nil {
		return err
	}

	for _, chapter := range selected {
		counter++

		s.Suffix = fmt.Sprintf(" [%d/%d] Getting pages of %s", counter, len(selected), style.Trim(40)(chapter.Name))
		_, err = src.PagesOf(chapter)
		if err != nil {
			return err
		}

		s.Suffix = fmt.Sprintf(" [%d/%d] Downloading %d pages", counter, len(selected), len(chapter.Pages))
		err = chapter.DownloadPages()
		if err != nil {
			return err
		}

		s.Suffix = fmt.Sprintf(" [%d/%d] Converting to %s", counter, len(selected), viper.GetString(config.FormatsUse))
		_, err = conv.Save(chapter)

		if viper.GetBool(config.HistorySaveOnDownload) {
			s.Suffix = fmt.Sprintf(" [%d/%d] Writing history", counter, len(selected))
			_ = history.Save(chapter)
		}
	}

	s.Stop()
	return nil
}
