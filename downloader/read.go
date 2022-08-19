package downloader

import (
	"fmt"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/converter"
	"github.com/metafates/mangal/history"
	"github.com/metafates/mangal/log"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/style"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/viper"
)

func Read(src source.Source, chapter *source.Chapter, progress func(string)) error {

	if viper.GetBool(config.ReaderReadInBrowser) {
		return open.Start(chapter.URL)
	}

	log.Info("downloading " + chapter.Name + " from " + chapter.Manga.Name + " for reading. Provider is " + src.ID())
	log.Info("getting pages of " + chapter.Name)
	progress("Getting pages")
	pages, err := src.PagesOf(chapter)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Info("downloading " + fmt.Sprintf("%d", len(pages)) + " pages")
	progress(fmt.Sprintf("Downloading %d pages", len(pages)))
	err = chapter.DownloadPages()
	if err != nil {
		log.Error(err)
		return err
	}

	log.Info("getting " + viper.GetString(config.FormatsUse) + " converter")
	conv, err := converter.Get(viper.GetString(config.FormatsUse))
	if err != nil {
		log.Error(err)
		return err
	}

	log.Info("converting " + viper.GetString(config.FormatsUse))
	progress(fmt.Sprintf(
		"Converting %d pages to %s %s",
		len(pages),
		style.Yellow(viper.GetString(config.FormatsUse)),
		style.Faint(chapter.SizeHuman())),
	)
	path, err := conv.SaveTemp(chapter)
	if err != nil {
		log.Error(err)
		return err
	}

	err = openRead(path, progress)
	if err != nil {
		return err
	}

	if viper.GetBool(config.HistorySaveOnRead) {
		go func() {
			err := history.Save(chapter)
			if err != nil {
				log.Warn(err)
			} else {
				log.Info("history saved")
			}
		}()
	}

	progress("Done")
	return nil
}

func openRead(path string, progress func(string)) error {
	var (
		reader string
		err    error
	)

	switch viper.GetString(config.FormatsUse) {
	case "pdf":
		reader = viper.GetString(config.ReaderPDF)
	case "cbz":
		reader = viper.GetString(config.ReaderCBZ)
	case "zip":
		reader = viper.GetString(config.ReaderZIP)
	case "plain":
		reader = viper.GetString(config.RaderPlain)
	}

	if reader != "" {
		log.Info("opening with " + reader)
		progress(fmt.Sprintf("Opening %s", reader))
		err = open.RunWith(reader, path)
		if err != nil {
			log.Error(err)
			return err
		}
		log.Info("opened without errors")
	} else {
		log.Info("no reader specified. opening with default")
		progress("Opening")
		err = open.Run(path)
		if err != nil {
			log.Error(err)
			return err
		}
		log.Info("opened without errors")
	}

	return nil
}
