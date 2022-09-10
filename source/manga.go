package source

import (
	"github.com/metafates/mangal/anilist"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/log"
	"github.com/metafates/mangal/util"
	"github.com/metafates/mangal/where"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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
	Chapters []*Chapter
	Source   Source `json:"-"`
	Metadata struct {
		Genres     []string
		Summary    string
		Author     string
		Cover      string
		Tags       []string
		Characters []string
		Status     string
		StartDate  struct {
			Year  int
			Month int
			Day   int
		}
	}
	cachedTempPath  string
	populated       bool
	coverDownloaded bool
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

func (m *Manga) DownloadCover(progress func(string)) error {
	if m.coverDownloaded {
		return nil
	}

	log.Info("Downloading cover for ", m.Name)
	progress("Downloading cover")

	m.coverDownloaded = true
	if m.Metadata.Cover == "" {
		log.Warn("No cover to download")
		return nil
	}

	path, err := m.Path(false)
	if err != nil {
		log.Error(err)
		return err
	}

	extension := ".jpg"
	if ext := filepath.Ext(m.Metadata.Cover); ext != "" {
		extension = ext
	}

	path = filepath.Join(path, "cover"+extension)

	exists, err := filesystem.Get().Exists(path)
	if err != nil {
		log.Error(err)
		return err
	}

	if exists {
		log.Warn("Cover already exists")
		return nil
	}

	resp, err := http.Get(m.Metadata.Cover)
	if err != nil {
		log.Error(err)
		return err
	}

	defer util.Ignore(resp.Body.Close)

	if resp.StatusCode != http.StatusOK {
		log.Error(err)
		return err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return err
	}

	return filesystem.Get().WriteFile(path, data, os.ModePerm)
}

func (m *Manga) PopulateMetadata() error {
	if m.populated {
		return nil
	}

	log.Infof("Populating metadata for %s", m.Name)

	manga, err := anilist.FindClosest(m.Name)
	if err != nil {
		log.Error(err)
		return err
	}

	m.Metadata.Genres = manga.Genres
	// replace <br> with newlines and remove other html tags
	m.Metadata.Summary = regexp.
		MustCompile("<.*?>").
		ReplaceAllString(
			strings.
				ReplaceAll(
					manga.Description,
					"<br>",
					"\n",
				),
			"",
		)

	var characters = make([]string, len(manga.Characters.Nodes))
	for i, character := range manga.Characters.Nodes {
		characters[i] = character.Name.Full
	}
	m.Metadata.Characters = characters

	var tags = make([]string, len(manga.Tags))
	for i, tag := range manga.Tags {
		tags[i] = tag.Name
	}
	m.Metadata.Tags = tags

	m.Metadata.Cover = manga.CoverImage.ExtraLarge
	m.Metadata.StartDate = struct {
		Year  int
		Month int
		Day   int
	}(manga.StartDate)

	m.Metadata.Status = strings.ReplaceAll(manga.Status, "_", " ")

	m.populated = true
	return nil
}
