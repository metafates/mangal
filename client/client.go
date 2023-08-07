package client

import (
	"context"
	"net/http"
	"time"

	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/anilist"
	"github.com/mangalorg/mangal/fs"
	"github.com/mangalorg/mangal/nametemplate"
)

func NewClient(ctx context.Context, loader libmangal.ProviderLoader) (*libmangal.Client, error) {
	HTTPClient := &http.Client{
		Timeout: time.Minute,
	}

	options := libmangal.DefaultClientOptions()
	options.FS = fs.Afero
	options.Anilist = anilist.Client
	options.HTTPClient = HTTPClient
	options.MangaNameTemplate = nametemplate.Manga
	options.VolumeNameTemplate = nametemplate.Volume
	options.ChapterNameTemplate = nametemplate.Chapter

	return libmangal.NewClient(ctx, loader, options)
}
