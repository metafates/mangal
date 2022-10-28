package source

import (
	"bytes"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/style"
	"github.com/metafates/mangal/util"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/spf13/viper"
	"html"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"
)

// Chapter is a struct that represents a chapter of a manga.
type Chapter struct {
	// Name of the chapter
	Name string
	// URL of the chapter
	URL string
	// Index of the chapter in the manga.
	Index uint16
	// ID of the chapter in the source.
	ID string
	// Volume which the chapter belongs to.
	Volume string
	// Manga that the chapter belongs to.
	Manga *Manga `json:"-"`
	// Pages of the chapter.
	Pages []*Page

	isDownloaded mo.Option[bool]
	size         uint64
}

func (c *Chapter) String() string {
	return c.Name
}

// DownloadPages downloads the Pages contents of the Chapter.
// Pages needs to be set before calling this function.
func (c *Chapter) DownloadPages(temp bool, progress func(string)) (err error) {
	c.size = 0
	status := func() string {
		return fmt.Sprintf(
			"Downloading %s %s",
			util.Quantify(len(c.Pages), "page", "pages"),
			style.Faint(c.SizeHuman()),
		)
	}

	progress(status())
	wg := sync.WaitGroup{}
	wg.Add(len(c.Pages))

	for _, page := range c.Pages {
		d := func(page *Page) {
			defer wg.Done()

			// if at any point, an error is encountered, stop downloading other pages
			if err != nil {
				return
			}

			err = page.Download()
			c.size += page.Size
			progress(status())
		}

		if viper.GetBool(constant.DownloaderAsync) {
			go d(page)
		} else {
			d(page)
		}
	}

	wg.Wait()
	c.isDownloaded = mo.Some(!temp && err == nil)
	return
}

// formattedName of the chapter according to the template in the config.
func (c *Chapter) formattedName() (name string) {
	name = viper.GetString(constant.DownloaderChapterNameTemplate)

	var sourceName string
	if c.Source() != nil {
		sourceName = c.Source().Name()
	}

	for variable, value := range map[string]string{
		"manga":          c.Manga.Name,
		"chapter":        c.Name,
		"index":          fmt.Sprintf("%d", c.Index),
		"padded-index":   fmt.Sprintf("%04d", c.Index),
		"chapters-count": fmt.Sprintf("%d", len(c.Manga.Chapters)),
		"volume":         c.Volume,
		"source":         sourceName,
	} {
		name = strings.ReplaceAll(name, fmt.Sprintf("{%s}", variable), value)
	}

	return
}

// SizeHuman is the same as Size but returns a human-readable string.
func (c *Chapter) SizeHuman() string {
	return humanize.Bytes(c.size)
}

func (c *Chapter) Filename() (filename string) {
	filename = util.SanitizeFilename(c.formattedName())

	// plain format assumes that chapter is a directory with images
	// rather than a single file. So no need to add extension to it
	if f := viper.GetString(constant.FormatsUse); f != constant.Plain {
		return filename + "." + f
	}

	return
}

func (c *Chapter) IsDownloaded() bool {
	if c.isDownloaded.IsPresent() {
		return c.isDownloaded.MustGet()
	}

	path, _ := c.path(c.Manga.peekPath(), false)
	exists, _ := filesystem.Api().Exists(path)
	c.isDownloaded = mo.Some(exists)
	return exists
}

func (c *Chapter) path(relativeTo string, createVolumeDir bool) (path string, err error) {
	if createVolumeDir {
		path = filepath.Join(path, util.SanitizeFilename(c.Volume))
		err = filesystem.Api().MkdirAll(path, os.ModePerm)
		if err != nil {
			return
		}
	}

	path = filepath.Join(relativeTo, c.Filename())
	return
}

func (c *Chapter) Path(temp bool) (path string, err error) {
	var manga string
	manga, err = c.Manga.Path(temp)
	if err != nil {
		return
	}

	return c.path(manga, c.Volume != "" && viper.GetBool(constant.DownloaderCreateVolumeDir))
}

func (c *Chapter) Source() Source {
	return c.Manga.Source
}

func (c *Chapter) ComicInfoXML() *bytes.Buffer {
	// language=gotemplate
	t := `
<ComicInfo xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	<Title>{{ escape .Name }}</Title>
  	<Series>{{ escape .Manga.Name }}</Series>
	<Number>{{ .Index }}</Number>
	<Web>{{ .URL }}</Web>
	<Genre>{{ join .Manga.Metadata.Genres "," }}</Genre>
	<PageCount>{{ len .Pages }}</PageCount>
	<Summary>{{ escape .Manga.Metadata.Summary }}</Summary>
	<Count>{{ len .Manga.Chapters }}</Count>
	<Writer>{{ escape .Manga.Metadata.Author }}</Writer>
	<Characters>{{ join .Manga.Metadata.Characters "," }}</Characters>
	{{ makeDate }}
	<Tags>{{ join .Manga.Metadata.Tags "," }}</Tags>
	<Notes>Downloaded with Mangal. https://github.com/metafates/mangal</Notes>
  	<Manga>YesAndRightToLeft</Manga>
</ComicInfo>`

	funcs := template.FuncMap{
		"join":   strings.Join,
		"escape": html.EscapeString,
		"geq":    func(a, b int) bool { return a >= b },
		"makeDate": func() string {
			if !viper.GetBool(constant.MetadataComicInfoXMLAddDate) {
				return ""
			}

			var (
				year  = lo.Tuple2[int, string]{0, "Year"}
				month = lo.Tuple2[int, string]{0, "Month"}
				day   = lo.Tuple2[int, string]{0, "Day"}
			)

			if viper.GetBool(constant.MetadataComicInfoXMLAlternativeDate) {
				// use current date (download date)
				now := time.Now()
				year.A = now.Year()
				month.A = int(now.Month())
				day.A = now.Day()
			} else {
				year.A = c.Manga.Metadata.StartDate.Year
				month.A = c.Manga.Metadata.StartDate.Month
				day.A = c.Manga.Metadata.StartDate.Day
			}

			sb := strings.Builder{}
			for _, t := range []lo.Tuple2[int, string]{year, month, day} {
				if t.A <= 0 {
					continue
				}

				sb.WriteString(fmt.Sprintf("<%[1]s>%d</%[1]s>\n", t.B, t.A))
			}

			return sb.String()
		},
	}

	parsed := lo.Must(template.New("ComicInfo").Funcs(funcs).Parse(t))
	buf := bytes.NewBufferString("")
	lo.Must0(parsed.Execute(buf, c))

	return buf
}
