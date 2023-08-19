package client

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/afs"
	"github.com/mangalorg/mangal/anilist"
	"github.com/mangalorg/mangal/nametemplate"
	"github.com/zyedidia/generic/queue"
)

var (
	clients = queue.New[*libmangal.Client]()
	m       sync.Mutex
)

func CloseAll() error {
	m.Lock()
	defer m.Unlock()

	for !clients.Empty() {
		client := clients.Peek()
		if err := client.Close(); err != nil {
			return err
		}

		clients.Dequeue()
	}

	return nil
}

func NewClient(ctx context.Context, loader libmangal.ProviderLoader) (*libmangal.Client, error) {
	HTTPClient := &http.Client{
		Timeout: time.Minute,
	}

	options := libmangal.DefaultClientOptions()
	options.FS = afs.Afero
	options.Anilist = anilist.Client
	options.HTTPClient = HTTPClient
	options.MangaNameTemplate = nametemplate.Manga
	options.VolumeNameTemplate = nametemplate.Volume
	options.ChapterNameTemplate = nametemplate.Chapter

	client, err := libmangal.NewClient(ctx, loader, options)
	if err != nil {
		return nil, err
	}

	clients.Enqueue(client)
	return client, nil
}
