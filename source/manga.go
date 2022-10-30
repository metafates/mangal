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
	"github.com/samber/mo"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type date struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

// Manga is a manga from a source.
type Manga struct {
	// Name of the manga
	Name string `json:"name"`
	// URL of the manga
	URL string `json:"url"`
	// Index of the manga in the source.
	Index uint16 `json:"index"`
	// ID of manga in the source.
	ID string `json:"id"`
	// Chapters of the manga
	Chapters []*Chapter `json:"chapters"`
	// Source that the manga belongs to.
	Source Source `json:"-"`
	// Anilist is the closest anilist match
	Anilist  mo.Option[*anilist.Manga] `json:"-"`
	Metadata struct {
		Genres  []string `json:"genres"`
		Summary string   `json:"summary"`
		Staff   struct {
			Story       []string `json:"story"`
			Art         []string `json:"art"`
			Translation []string `json:"translation"`
			Lettering   []string `json:"lettering"`
		} `json:"staff"`
		Cover struct {
			ExtraLarge string `json:"extraLarge"`
			Large      string `json:"large"`
			Medium     string `json:"medium"`
			Color      string `json:"color"`
		} `json:"cover"`
		Tags       []string `json:"tags"`
		Characters []string `json:"characters"`
		Status     string   `json:"status"`
		StartDate  date     `json:"startDate"`
		EndDate    date     `json:"endDate"`
		Synonyms   []string `json:"synonyms"`
		URLs       []string `json:"urls"`
	} `json:"metadata"`
	cachedTempPath  string
	populated       bool
	coverDownloaded bool
}

func (m *Manga) String() string {
	return m.Name
}

func (m *Manga) Dirname() string {
	return util.SanitizeFilename(m.Name)
}

func (m *Manga) peekPath() string {
	path := where.Downloads()

	if viper.GetBool(constant.DownloaderCreateMangaDir) {
		path = filepath.Join(path, m.Dirname())
	}

	return path
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

	path = m.peekPath()
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

	if m.Metadata.Cover.ExtraLarge == "" {
		log.Warn("No cover to download")
		return nil
	}

	path, err := m.Path(false)
	if err != nil {
		log.Error(err)
		return err
	}

	var extension string
	if extension = filepath.Ext(m.Metadata.Cover.ExtraLarge); extension == "" {
		extension = ".jpg"
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

	resp, err := http.Get(m.Metadata.Cover.ExtraLarge)
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

func (m *Manga) BindWithAnilist() error {
	if m.Anilist.IsPresent() {
		return nil
	}

	log.Infof("binding %s with anilist", m.Name)

	manga, err := anilist.FindClosest(m.Name)
	if err != nil {
		log.Error(err)
		return err
	}

	m.Anilist = mo.Some(manga)
	return nil
}

func (m *Manga) PopulateMetadata(progress func(string)) error {
	if m.populated {
		return nil
	}
	m.populated = true

	progress("Fetching metadata from anilist")
	log.Infof("Populating metadata for %s", m.Name)
	if err := m.BindWithAnilist(); err != nil {
		progress("Failed to fetch metadata")
		return err
	}

	manga, ok := m.Anilist.Get()
	if !ok || manga == nil {
		return fmt.Errorf("manga '%s' not found on Anilist", m.Name)
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

	var tags = make([]string, 0)
	for _, tag := range manga.Tags {
		if tag.Rank >= viper.GetInt(constant.MetadataComicInfoXMLTagRelevanceThreshold) {
			tags = append(tags, tag.Name)
		}
	}
	m.Metadata.Tags = tags

	m.Metadata.Cover.ExtraLarge = manga.CoverImage.ExtraLarge
	m.Metadata.Cover.Medium = manga.CoverImage.Medium
	m.Metadata.Cover.Color = manga.CoverImage.Color

	m.Metadata.StartDate = date(manga.StartDate)
	m.Metadata.EndDate = date(manga.EndDate)

	m.Metadata.Status = strings.ReplaceAll(manga.Status, "_", " ")
	m.Metadata.Synonyms = manga.Synonyms

	m.Metadata.Staff.Story = make([]string, 0)
	m.Metadata.Staff.Art = make([]string, 0)
	m.Metadata.Staff.Translation = make([]string, 0)
	m.Metadata.Staff.Lettering = make([]string, 0)

	for _, staff := range manga.Staff.Edges {
		role := strings.ToLower(staff.Role)
		switch {
		case strings.Contains(role, "story"):
			m.Metadata.Staff.Story = append(m.Metadata.Staff.Story, staff.Node.Name.Full)
		case strings.Contains(role, "art"):
			m.Metadata.Staff.Art = append(m.Metadata.Staff.Art, staff.Node.Name.Full)
		case strings.Contains(role, "translator"):
			m.Metadata.Staff.Translation = append(m.Metadata.Staff.Translation, staff.Node.Name.Full)
		case strings.Contains(role, "lettering"):
			m.Metadata.Staff.Lettering = append(m.Metadata.Staff.Lettering, staff.Node.Name.Full)
		}
	}

	// Anilist & Myanimelist + external
	urls := make([]string, 2+len(manga.External))
	urls[0] = manga.SiteURL
	for i, e := range manga.External {
		urls[i+1] = e.URL
	}

	urls = lo.Filter(urls, func(url string, _ int) bool {
		return url != ""
	})

	urls = append(urls, fmt.Sprintf("https://myanimelist.net/manga/%d", manga.IDMal))
	m.Metadata.URLs = urls

	return nil
}

func (m *Manga) SeriesJSON() *bytes.Buffer {
	var status string
	switch m.Metadata.Status {
	case "FINISHED":
		status = "Ended"
	case "RELEASING":
		status = "Continuing"
	default:
		status = "Unknown"
	}

	var publisher string
	if len(m.Metadata.Staff.Story) > 0 {
		publisher = m.Metadata.Staff.Story[0]
	}

	seriesJSON := &SeriesJSON{}
	seriesJSON.Metadata.Type = "comicSeries"
	seriesJSON.Metadata.Name = m.Name
	seriesJSON.Metadata.DescriptionFormatted = m.Metadata.Summary
	seriesJSON.Metadata.Status = status
	seriesJSON.Metadata.Year = m.Metadata.StartDate.Year
	seriesJSON.Metadata.ComicImage = m.Metadata.Cover.ExtraLarge
	seriesJSON.Metadata.Publisher = publisher
	seriesJSON.Metadata.BookType = "Print"
	seriesJSON.Metadata.TotalIssues = len(m.Chapters)
	seriesJSON.Metadata.PublicationRun = fmt.Sprintf("%d %d - %d %d", m.Metadata.StartDate.Month, m.Metadata.StartDate.Year, m.Metadata.EndDate.Month, m.Metadata.EndDate.Year)

	var buf bytes.Buffer
	lo.Must0(json.NewEncoder(&buf).Encode(seriesJSON))
	return &buf
}
