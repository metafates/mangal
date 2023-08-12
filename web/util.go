package web

import (
	"context"
	"fmt"
	"strconv"

	"github.com/mangalorg/libmangal"
)

func searchMangas(ctx context.Context, client *libmangal.Client, query string) ([]libmangal.Manga, error) {
	return client.SearchMangas(ctx, query)
}

func mangaVolumes(ctx context.Context, client *libmangal.Client, query, mangaID string) ([]libmangal.Volume, error) {
	mangas, err := searchMangas(ctx, client, query)
	if err != nil {
		return nil, err
	}

	for _, manga := range mangas {
		if manga.Info().ID == mangaID {
			return client.MangaVolumes(ctx, manga)
		}
	}

	return nil, fmt.Errorf("manga %q not found", mangaID)
}

func volumeChapters(ctx context.Context, client *libmangal.Client, query, mangaID string, volumeNumber int) ([]libmangal.Chapter, error) {
	volumes, err := mangaVolumes(ctx, client, query, mangaID)
	if err != nil {
		return nil, err
	}

	for _, volume := range volumes {
		if volume.Info().Number == volumeNumber {
			return client.VolumeChapters(ctx, volume)
		}
	}

	return nil, fmt.Errorf("volume %q not found", strconv.Itoa(volumeNumber))
}
