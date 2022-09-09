package source

import (
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/util"
	"github.com/metafates/mangal/where"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

// Manga is a manga from a source.
type Manga struct {
	// Name of the manga
	Name string
	// URL of the manga
	URL string
	// Index of the manga in the source.
	Index uint16
	// ID of manga in the source.
	ID string
	// Chapters of the manga
	Chapters       []*Chapter
	Source         Source `json:"-"`
	Genres         []string
	Summary        string
	Author         string
	cachedTempPath string
}

func (m *Manga) String() string {
	return m.Name
}

func (m *Manga) Filename() string {
	return util.SanitizeFilename(m.Name)
}

func (m *Manga) Path(temp bool) (path string, err error) {
	if temp {
		if path = m.cachedTempPath; path != "" {
			return
		}

		path, err = filesystem.Get().TempDir("", constant.TempPrefix)

		m.cachedTempPath = path
		return
	}

	path = where.Downloads()

	if viper.GetBool(constant.DownloaderCreateMangaDir) {
		path = filepath.Join(path, m.Filename())
	}

	_ = filesystem.Get().MkdirAll(path, os.ModePerm)
	return
}
