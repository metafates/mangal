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
	"github.com/spf13/viper"
	"html"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
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

	size uint64
}

func (c *Chapter) String() string {
	return c.Name
}

// DownloadPages downloads the Pages contents of the Chapter.
// Pages needs to be set before calling this function.
func (c *Chapter) DownloadPages(progress func(string)) (err error) {
	c.size = 0
	status := func() string {
		return fmt.Sprintf(
			"Downloading %s %s",
			util.Quantity(len(c.Pages), "page"),
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

func (c *Chapter) Path(temp bool) (path string, err error) {
	path, err = c.Manga.Path(temp)
	if err != nil {
		return
	}

	if c.Volume != "" && viper.GetBool(constant.DownloaderCreateVolumeDir) {
		path = filepath.Join(path, util.SanitizeFilename(c.Volume))
		err = filesystem.Api().MkdirAll(path, os.ModePerm)
		if err != nil {
			return
		}
	}

	path = filepath.Join(path, c.Filename())
	return
}

func (c *Chapter) Source() Source {
	return c.Manga.Source
}

func (c *Chapter) ComicInfoXML() *bytes.Buffer {
	// language=gotemplate
	t := `
<ComicInfo xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	<Title>{{ escape .Key }}</Title>
  	<Series>{{ escape .Manga.Key }}</Series>
	<Number>{{ .Index }}</Number>
	<Web>{{ .URL }}</Web>
	<Genre>{{ join .Manga.Metadata.Genres "," }}</Genre>
	<PageCount>{{ len .Pages }}</PageCount>
	<Summary>{{ escape .Manga.Metadata.Summary }}</Summary>
	<Count>{{ len .Manga.Chapters }}</Count>
	<Writer>{{ .Manga.Metadata.Author }}</Writer>
	<Characters>{{ join .Manga.Metadata.Characters "," }}</Characters>
	<Year>{{ .Manga.Metadata.StartDate.Year }}</Year>
	<Month>{{ .Manga.Metadata.StartDate.Month }}</Month>
	<Day>{{ .Manga.Metadata.StartDate.Day }}</Day>
	<Tags>{{ join .Manga.Metadata.Tags "," }}</Tags>
	<Notes>Downloaded with Mangal. https://github.com/metafates/mangal</Notes>
  	<Manga>YesAndRightToLeft</Manga>
</ComicInfo>`

	funcs := template.FuncMap{
		"join":   strings.Join,
		"escape": html.EscapeString,
	}

	parsed := lo.Must(template.New("ComicInfo").Funcs(funcs).Parse(t))
	buf := bytes.NewBufferString("")
	lo.Must0(parsed.Execute(buf, c))

	return buf
}
