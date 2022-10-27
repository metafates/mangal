package downloader

import (
	"fmt"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/converter"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/history"
	"github.com/metafates/mangal/log"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/style"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

// Download the chapter using given source.
func Download(chapter *source.Chapter, progress func(string)) (string, error) {
	log.Info("downloading " + chapter.Name)

	log.Info("checking if chapter is already downloaded")
	path, err := chapter.Path(false)
	if err != nil {
		log.Error(err)
		return "", err
	}

	exists, err := filesystem.Api().Exists(path)
	if err != nil {
		log.Error(err)
		return "", err
	}

	if viper.GetBool(constant.DownloaderRedownloadExisting) {
		log.Info("chapter already downloaded, deleting and redownloading")
		err = filesystem.Api().Remove(path)
		if err != nil {
			log.Warn(err)
		}
	} else {
		if exists {
			log.Info("chapter already downloaded, skipping")
			return path, nil
		}
	}

	progress("Getting pages")
	pages, err := chapter.Source().PagesOf(chapter)
	if err != nil {
		log.Error(err)
		return "", err
	}
	log.Info("found " + fmt.Sprintf("%d", len(pages)) + " pages")

	err = chapter.DownloadPages(false, progress)
	if err != nil {
		log.Error(err)
		return "", err
	}

	if viper.GetBool(constant.MetadataFetchAnilist) {
		err := chapter.Manga.PopulateMetadata(progress)
		if err != nil {
			log.Warn(err)
		}
	}

	if viper.GetBool(constant.MetadataSeriesJSON) {
		path, err := chapter.Manga.Path(false)
		if err != nil {
			log.Warn(err)
		} else {
			path = filepath.Join(path, "series.json")
			progress("Generating series.json")
			err = filesystem.Api().WriteFile(path, chapter.Manga.SeriesJSON().Bytes(), os.ModePerm)
			if err != nil {
				log.Warn(err)
			}
		}
	}

	if viper.GetBool(constant.DownloaderDownloadCover) {
		_ = chapter.Manga.DownloadCover(progress)
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
	path, err = conv.Save(chapter)
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
