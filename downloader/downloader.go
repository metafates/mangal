package downloader

import (
	"github.com/metafates/mangal/source"
	"sync"
)

// Populate will download the given pages and set contents to each
func Populate(pages []*source.Page) error {
	wg := sync.WaitGroup{}
	wg.Add(len(pages))

	var err error
	for _, page := range pages {
		go func(page *source.Page) {
			defer wg.Done()

			// if at any point, an error is encountered, stop downloading other pages
			if err != nil {
				return
			}

			err = page.Download()
		}(page)
	}

	wg.Wait()
	return err
}
