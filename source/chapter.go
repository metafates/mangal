package source

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/util"
	"github.com/spf13/viper"
	"strings"
	"sync"
)

type Chapter struct {
	Name     string
	URL      string
	Index    uint16
	SourceID string
	ID       string
	Manga    *Manga
	Pages    []*Page
}

func (c *Chapter) String() string {
	return c.Name
}

// DownloadPages downloads the Pages contents of the Chapter.
// Pages needs to be set before calling this function.
func (c *Chapter) DownloadPages() error {
	wg := sync.WaitGroup{}
	wg.Add(len(c.Pages))

	var err error
	for _, page := range c.Pages {
		d := func(page *Page) {
			defer wg.Done()

			// if at any point, an error is encountered, stop downloading other pages
			if err != nil {
				return
			}

			err = page.Download()
		}

		if viper.GetBool(config.DownloaderAsync) {
			go d(page)
		} else {
			d(page)
		}
	}

	wg.Wait()
	return err
}

func (c *Chapter) FormattedName() (name string) {
	template := viper.GetString(config.DownloaderChapterNameTemplate)
	name = strings.ReplaceAll(template, "{manga}", c.Manga.Name)
	name = strings.ReplaceAll(name, "{chapter}", c.Name)
	name = strings.ReplaceAll(name, "{index}", fmt.Sprintf("%d", c.Index))
	name = strings.ReplaceAll(name, "{padded-index}", util.PadZero(fmt.Sprintf("%d", c.Index), 4))

	return
}

func (c *Chapter) Size() uint64 {
	var n uint64

	for _, page := range c.Pages {
		n += page.Size
	}

	return n
}

func (c *Chapter) SizeHuman() string {
	if size := c.Size(); size == 0 {
		return "Unknown size"
	} else {
		return humanize.Bytes(size)
	}
}
