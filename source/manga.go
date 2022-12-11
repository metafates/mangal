package source

import (
	"fmt"
	"github.com/metafates/mangal/anilist"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/key"
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
	Name string `json:"name" jsonschema:"description=Name of the manga"`
	// URL of the manga
	URL string `json:"url" jsonschema:"description=URL of the manga"`
	// Index of the manga in the source.
	Index uint16 `json:"index" jsonschema:"description=Index of the manga in the source"`
	// ID of manga in the source.
	ID string `json:"id" jsonschema:"description=ID of manga in the source"`
	// Chapters of the manga
	Chapters []*Chapter `json:"chapters" jsonschema:"description=Chapters of the manga"`
	// Source that the manga belongs to.
	Source Source `json:"-"`
	// Anilist is the closest anilist match
	Anilist  mo.Option[*anilist.Manga] `json:"-"`
	Metadata struct {
		// Genres of the manga
		Genres []string `json:"genres" jsonschema:"description=Genres of the manga"`
		// Summary in the plain text with newlines
		Summary string `json:"summary" jsonschema:"description=Summary in the plain text with newlines"`
		// Staff that worked on the manga
		Staff struct {
			// Story authors
			Story []string `json:"story" jsonschema:"description=Story authors"`
			// Art authors
			Art []string `json:"art" jsonschema:"description=Art authors"`
			// Translation group
			Translation []string `json:"translation" jsonschema:"description=Translation group"`
			// Lettering group
			Lettering []string `json:"lettering" jsonschema:"description=Lettering group"`
		} `json:"staff" jsonschema:"description=Staff that worked on the manga"`
		// Cover images of the manga
		Cover struct {
			// ExtraLarge is the largest cover image. If not available, Large will be used.
			ExtraLarge string `json:"extraLarge" jsonschema:"description=ExtraLarge is the largest cover image. If not available, Large will be used."`
			// Large is the second-largest cover image.
			Large string `json:"large" jsonschema:"description=Large is the second-largest cover image."`
			// Medium cover image. The smallest one.
			Medium string `json:"medium" jsonschema:"description=Medium cover image. The smallest one."`
			// Color average color of the cover image.
			Color string `json:"color" jsonschema:"description=Color average color of the cover image."`
		} `json:"cover" jsonschema:"description=Cover images of the manga"`
		// BannerImage is the banner image of the manga.
		BannerImage string `json:"bannerImage" jsonschema:"description=BannerImage is the banner image of the manga."`
		// Tags of the manga
		Tags []string `json:"tags" jsonschema:"description=Tags of the manga"`
		// Characters of the manga
		Characters []string `json:"characters" jsonschema:"description=Characters of the manga"`
		// Status of the manga
		Status string `json:"status" jsonschema:"enum=FINISHED,enum=RELEASING,enum=NOT_YET_RELEASED,enum=CANCELLED,enum=HIATUS"`
		// StartDate is the date when the manga started.
		StartDate date `json:"startDate" jsonschema:"description=StartDate is the date when the manga started."`
		// EndDate is the date when the manga ended.
		EndDate date `json:"endDate" jsonschema:"description=EndDate is the date when the manga ended."`
		// Synonyms other names of the manga.
		Synonyms []string `json:"synonyms" jsonschema:"description=Synonyms other names of the manga."`
		// Chapters is the amount of chapters the manga will have when completed.
		Chapters int `json:"chapters" jsonschema:"description=The amount of chapters the manga will have when completed."`
		// URLs external URLs of the manga.
		URLs []string `json:"urls" jsonschema:"description=External URLs of the manga."`
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

	if viper.GetBool(key.DownloaderCreateMangaDir) {
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

func (m *Manga) GetCover() (string, error) {
	var covers = []string{
		m.Metadata.Cover.ExtraLarge,
		m.Metadata.Cover.Large,
		m.Metadata.Cover.Medium,
	}

	for _, cover := range covers {
		if cover != "" {
			return cover, nil
		}
	}

	return "", fmt.Errorf("no cover found")
}

func (m *Manga) DownloadCover(overwrite bool, path string, progress func(string)) error {
	if m.coverDownloaded {
		return nil
	}
	m.coverDownloaded = true

	log.Info("Downloading cover for ", m.Name)
	progress("Downloading cover")

	cover, err := m.GetCover()
	if err != nil {
		log.Warn(err)
		return nil
	}

	var extension string
	if extension = filepath.Ext(cover); extension == "" {
		extension = ".jpg"
	}

	path = filepath.Join(path, "cover"+extension)

	if !overwrite {
		exists, err := filesystem.Api().Exists(path)
		if err != nil {
			log.Error(err)
			return err
		}

		if exists {
			log.Warn("Cover already exists")
			return nil
		}
	}

	resp, err := http.Get(cover)
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

	err = filesystem.Api().WriteFile(path, data, os.ModePerm)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Info("Cover downloaded")
	return nil
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
		if tag.Rank >= viper.GetInt(key.MetadataComicInfoXMLTagRelevanceThreshold) {
			tags = append(tags, tag.Name)
		}
	}
	m.Metadata.Tags = tags

	m.Metadata.Cover.ExtraLarge = manga.CoverImage.ExtraLarge
	m.Metadata.Cover.Large = manga.CoverImage.Large
	m.Metadata.Cover.Medium = manga.CoverImage.Medium
	m.Metadata.Cover.Color = manga.CoverImage.Color

	m.Metadata.BannerImage = manga.BannerImage

	m.Metadata.StartDate = date(manga.StartDate)
	m.Metadata.EndDate = date(manga.EndDate)

	m.Metadata.Status = strings.ReplaceAll(manga.Status, "_", " ")
	m.Metadata.Synonyms = manga.Synonyms

	m.Metadata.Staff.Story = make([]string, 0)
	m.Metadata.Staff.Art = make([]string, 0)
	m.Metadata.Staff.Translation = make([]string, 0)
	m.Metadata.Staff.Lettering = make([]string, 0)

	m.Metadata.Chapters = manga.Chapters

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

func (m *Manga) SeriesJSON() *SeriesJSON {
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
	seriesJSON.Metadata.DescriptionText = m.Metadata.Summary
	seriesJSON.Metadata.Status = status
	seriesJSON.Metadata.Year = m.Metadata.StartDate.Year
	seriesJSON.Metadata.ComicImage = m.Metadata.Cover.ExtraLarge
	seriesJSON.Metadata.Publisher = publisher
	seriesJSON.Metadata.BookType = "Print"
	seriesJSON.Metadata.TotalIssues = m.Metadata.Chapters
	seriesJSON.Metadata.PublicationRun = fmt.Sprintf("%d %d - %d %d", m.Metadata.StartDate.Month, m.Metadata.StartDate.Year, m.Metadata.EndDate.Month, m.Metadata.EndDate.Year)

	return seriesJSON
}
