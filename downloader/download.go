package downloader

import (
	"fmt"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/converter"
	"github.com/metafates/mangal/history"
	"github.com/metafates/mangal/log"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/style"
	"github.com/spf13/viper"
)

// Download the chapter using given source.
func Download(chapter *source.Chapter, progress func(string)) (string, error) {
	log.Info("downloading " + chapter.Name)
	progress("Getting pages")
	pages, err := chapter.Source().PagesOf(chapter)
	if err != nil {
		log.Error(err)
		return "", err
	}
	log.Info("found " + fmt.Sprintf("%d", len(pages)) + " pages")

	err = chapter.DownloadPages(progress)
	if err != nil {
		log.Error(err)
		return "", err
	}

	log.Info("getting " + viper.GetString(constant.FormatsUse) + " converter")
	progress(fmt.Sprintf(
		"Converting %d pages to %s %s",
		len(pages),
		style.Yellow(viper.GetString(constant.FormatsUse)),
		style.Faint(chapter.SizeHuman())),
	)
	conv, err := converter.Get(viper.GetString(constant.FormatsUse))
	if err != nil {
		log.Error(err)
		return "", err
	}

	log.Info("converting " + viper.GetString(constant.FormatsUse))
	path, err := conv.Save(chapter)
	if err != nil {
		log.Error(err)
		return "", err
	}

	if viper.GetBool(constant.HistorySaveOnDownload) {
		go func() {
			err = history.Save(chapter)
			if err != nil {
				log.Warn(err)
			} else {
				log.Info("history saved")
			}
		}()
	}

	log.Info("downloaded without errors")
	progress("Downloaded")
	return path, nil
}
