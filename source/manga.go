package source

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/metafates/mangal/anilist"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/log"
	"github.com/metafates/mangal/util"
	"github.com/metafates/mangal/where"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type date struct {
	Year  int
	Month int
	Day   int
}

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
		StartDate  date
		EndDate    date
		Synonyms   []string
		URLs       []string
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

		path = where.Temp()
		m.cachedTempPath = path
		return
	}

	path = where.Downloads()

	if viper.GetBool(constant.DownloaderCreateMangaDir) {
		path = filepath.Join(path, m.Filename())
	}

	_ = filesystem.Api().MkdirAll(path, os.ModePerm)
	return
}

func (m *Manga) DownloadCover(progress func(string)) error {
	if m.coverDownloaded {
		return nil
	}
	m.coverDownloaded = true

	log.Info("Downloading cover for ", m.Name)
	progress("Downloading cover")

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

	exists, err := filesystem.Api().Exists(path)
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

	return filesystem.Api().WriteFile(path, data, os.ModePerm)
}

func (m *Manga) PopulateMetadata(progress func(string)) error {
	if m.populated {
		return nil
	}
	m.populated = true

	progress("Fetching metadata from anilist")
	log.Infof("Populating metadata for %s", m.Name)

	manga, err := anilist.FindClosest(m.Name)
	if err != nil {
		log.Error(err)
		progress("Failed to fetch metadata")
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
	m.Metadata.StartDate = date(manga.StartDate)
	m.Metadata.EndDate = date(manga.EndDate)

	m.Metadata.Status = strings.ReplaceAll(manga.Status, "_", " ")
	m.Metadata.Synonyms = manga.Synonyms

	urls := []string{manga.URL}
	if manga.SiteURL != "" {
		urls = append(urls, manga.SiteURL)
	}

	for _, e := range manga.External {
		if e.URL != "" {
			urls = append(urls, e.URL)
		}
	}

	urls = append(urls, fmt.Sprintf("https://myanimelist.net/manga/%d", manga.IDMal))
	m.Metadata.URLs = urls

	return nil
}

func (m *Manga) SeriesJSON() *bytes.Buffer {
	type metadata struct {
		Type                 string `json:"type"`
		Name                 string `json:"name"`
		DescriptionFormatted string `json:"description_formatted"`
		DescriptionText      string `json:"description_text"`
		Status               string `json:"status"`
		Year                 int    `json:"year"`
		ComicImage           string `json:"ComicImage"`
		Publisher            string `json:"publisher"`
		ComicID              int    `json:"comicId"`
		BookType             string `json:"booktype"`
		TotalIssues          int    `json:"total_issues"`
		PublicationRun       string `json:"publication_run"`
	}

	var status string
	switch m.Metadata.Status {
	case "FINISHED":
		status = "Ended"
	case "RELEASING":
		status = "Continuing"
	default:
		status = "Unknown"
	}

	seriesJSON := struct {
		Metadata metadata `json:"metadata"`
	}{
		Metadata: metadata{
			Type:                 "comicSeries",
			Name:                 m.Name,
			DescriptionFormatted: m.Metadata.Summary,
			Status:               status,
			Year:                 m.Metadata.StartDate.Year,
			ComicImage:           m.Metadata.Cover,
			Publisher:            m.Metadata.Author,
			BookType:             "Print",
			TotalIssues:          len(m.Chapters),
			PublicationRun:       fmt.Sprintf("%d %d - %d %d", m.Metadata.StartDate.Month, m.Metadata.StartDate.Year, m.Metadata.EndDate.Month, m.Metadata.EndDate.Year),
		},
	}

	var buf bytes.Buffer
	lo.Must0(json.NewEncoder(&buf).Encode(seriesJSON))
	return &buf
}
