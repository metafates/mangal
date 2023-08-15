package web

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/anilist"
	"github.com/mangalorg/mangal/client"
	"github.com/mangalorg/mangal/meta"
	"github.com/mangalorg/mangal/provider/manager"
	"github.com/mangalorg/mangal/web/api"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/bigcache"
	"github.com/philippgille/gokv/encoding"
	"github.com/samber/lo"
)

//go:embed all:ui/dist
var frontend embed.FS

var _ api.StrictServerInterface = (*Server)(nil)

type Server struct {
	imageCache  gokv.Store
	loaders     []libmangal.ProviderLoader
	loadersByID map[string]libmangal.ProviderLoader
}

func (s *Server) GetMangaPage(ctx context.Context, request api.GetMangaPageRequestObject) (api.GetMangaPageResponseObject, error) {
	loader, ok := s.loadersByID[request.Params.Provider]
	if !ok {
		return api.GetMangaPagedefaultJSONResponse{
			StatusCode: 404,
			Body: api.Error{
				Code:    404,
				Message: fmt.Sprintf("Provider %q not found", request.Params.Provider),
			},
		}, nil
	}

	c, err := client.NewClient(ctx, loader)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	mangas, err := searchMangas(ctx, c, request.Params.Query)
	if err != nil {
		return nil, err
	}

	manga, ok := lo.Find(mangas, func(manga libmangal.Manga) bool {
		return manga.Info().ID == request.Params.Manga
	})
	if !ok {
		return nil, fmt.Errorf("manga with id %s not found", request.Params.Manga)
	}

	volumes, err := c.MangaVolumes(ctx, manga)
	if err != nil {
		return nil, err
	}

	v := make([]api.VolumeWithChapters, len(volumes))
	for i, volume := range volumes {
		chapters, err := c.VolumeChapters(ctx, volume)
		if err != nil {
			return nil, err
		}

		c := make([]api.Chapter, len(chapters))
		for j, chapter := range chapters {
			info := chapter.Info()
			c[j] = api.Chapter{
				Number: info.Number,
				Title:  info.Title,
				Url:    &info.URL,
			}
		}

		v[i] = api.VolumeWithChapters{
			Chapters: c,
			Volume: api.Volume{
				Number: volume.Info().Number,
			},
		}
	}

	info := manga.Info()
	response := api.GetMangaPage200JSONResponse{
		Manga: api.Manga{
			Banner: &info.Banner,
			Cover:  &info.Cover,
			Id:     info.ID,
			Title:  info.Title,
			Url:    &info.URL,
		},
		Volumes: v,
	}

	anilistManga, found, _ := anilist.Client.FindClosestManga(ctx, info.Title)
	if found {
		response.AnilistManga = &api.AnilistManga{
			BannerImage: &anilistManga.BannerImage,
			CoverImage: api.CoverImage{
				Color:      anilistManga.CoverImage.Color,
				ExtraLarge: anilistManga.CoverImage.ExtraLarge,
				Large:      anilistManga.CoverImage.Large,
				Medium:     anilistManga.CoverImage.Medium,
			},
			Description: &anilistManga.Description,
		}
	}

	return response, nil
}

// GetChapter implements api.StrictServerInterface.
func (*Server) GetChapter(ctx context.Context, request api.GetChapterRequestObject) (api.GetChapterResponseObject, error) {
	panic("unimplemented")
}

// GetManga implements api.StrictServerInterface.
func (*Server) GetManga(ctx context.Context, request api.GetMangaRequestObject) (api.GetMangaResponseObject, error) {
	panic("unimplemented")
}

// GetFormats implements api.StrictServerInterface.
func (*Server) GetFormats(ctx context.Context, request api.GetFormatsRequestObject) (api.GetFormatsResponseObject, error) {
	return api.GetFormats200JSONResponse(
		lo.Map(libmangal.FormatValues(), func(format libmangal.Format, _ int) api.Format {
			f := api.Format{
				Name: format.String(),
			}

			if extension := format.Extension(); extension != "" {
				f.Extension = &extension
			}

			return f
		}),
	), nil
}

// GetProvider implements api.StrictServerInterface.
func (s *Server) GetProvider(ctx context.Context, request api.GetProviderRequestObject) (api.GetProviderResponseObject, error) {
	loader, ok := s.loadersByID[request.Params.Id]
	if !ok {
		return api.GetProvider404Response{}, nil
	}

	info := loader.Info()

	return api.GetProvider200JSONResponse{
		Description: &info.Description,
		Id:          info.ID,
		Name:        info.Name,
		Version:     info.Version,
	}, nil
}

// GetMangaVolumes implements api.StrictServerInterface.
func (s *Server) GetMangaVolumes(ctx context.Context, request api.GetMangaVolumesRequestObject) (api.GetMangaVolumesResponseObject, error) {
	loader, ok := s.loadersByID[request.Params.Provider]
	if !ok {
		return api.GetMangaVolumesdefaultJSONResponse{
			StatusCode: 404,
			Body: api.Error{
				Code:    404,
				Message: fmt.Sprintf("Provider %q not found", request.Params.Provider),
			},
		}, nil
	}

	c, err := client.NewClient(ctx, loader)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	volumes, err := mangaVolumes(ctx, c, request.Params.Query, request.Params.Manga)
	if err != nil {
		return nil, err
	}

	return api.GetMangaVolumes200JSONResponse(
		lo.Map(volumes, func(volume libmangal.Volume, _ int) api.Volume {
			return api.Volume{
				Number: volume.Info().Number,
			}
		}),
	), nil
}

// GetVolumeChapters implements api.StrictServerInterface.
func (s *Server) GetVolumeChapters(ctx context.Context, request api.GetVolumeChaptersRequestObject) (api.GetVolumeChaptersResponseObject, error) {
	loader, ok := s.loadersByID[request.Params.Provider]
	if !ok {
		return api.GetVolumeChaptersdefaultJSONResponse{
			StatusCode: 404,
			Body: api.Error{
				Code:    404,
				Message: fmt.Sprintf("Provider %q not found", request.Params.Provider),
			},
		}, nil
	}

	c, err := client.NewClient(ctx, loader)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	chapters, err := volumeChapters(ctx, c, request.Params.Query, request.Params.Manga, request.Params.Volume)
	if err != nil {
		return nil, err
	}

	return api.GetVolumeChapters200JSONResponse(
		lo.Map(chapters, func(chapter libmangal.Chapter, _ int) api.Chapter {
			info := chapter.Info()
			return api.Chapter{
				Number: info.Number,
				Title:  info.Title,
				Url:    &info.URL,
			}
		}),
	), nil
}

// GetImage implements api.StrictServerInterface.
func (s *Server) GetImage(ctx context.Context, request api.GetImageRequestObject) (api.GetImageResponseObject, error) {
	var image []byte
	found, err := s.imageCache.Get(request.Params.Url, &image)
	if err != nil {
		s.imageCache.Delete(request.Params.Url)
	}

	if found {
		return api.GetImage200ImagepngResponse{
			Body:          bytes.NewReader(image),
			ContentLength: int64(len(image)),
		}, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, request.Params.Url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
	if referer := request.Params.Referer; referer != nil {
		req.Header.Set("Referer", *referer)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	image, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	_ = s.imageCache.Set(request.Params.Url, image)

	return api.GetImage200ImagepngResponse{
		Body:          bytes.NewReader(image),
		ContentLength: int64(len(image)),
	}, nil
}

// GetMangalInfo implements api.StrictServerInterface.
func (*Server) GetMangalInfo(ctx context.Context, request api.GetMangalInfoRequestObject) (api.GetMangalInfoResponseObject, error) {
	return api.GetMangalInfo200JSONResponse{
		Version: meta.Version,
	}, nil
}

// SearchMangas implements api.StrictServerInterface.
func (s *Server) SearchMangas(ctx context.Context, request api.SearchMangasRequestObject) (api.SearchMangasResponseObject, error) {
	loader, ok := s.loadersByID[request.Params.Provider]

	if !ok {
		return api.SearchMangasdefaultJSONResponse{
			StatusCode: 404,
			Body: api.Error{
				Code:    404,
				Message: fmt.Sprintf("Provider %q not found", request.Params.Provider),
			},
		}, nil
	}

	c, err := client.NewClient(ctx, loader)
	if err != nil {
		return api.SearchMangasdefaultJSONResponse{
			StatusCode: 400,
			Body: api.Error{
				Message: err.Error(),
			},
		}, nil
	}
	defer c.Close()

	mangas, err := searchMangas(ctx, c, request.Params.Query)
	if err != nil {
		return api.SearchMangasdefaultJSONResponse{
			StatusCode: 400,
			Body: api.Error{
				Message: err.Error(),
			},
		}, nil
	}

	return api.SearchMangas200JSONResponse(lo.Map(mangas, func(manga libmangal.Manga, _ int) api.Manga {
		info := manga.Info()
		return api.Manga{
			Banner: &info.Banner,
			Cover:  &info.Cover,
			Id:     info.ID,
			Title:  info.Title,
			Url:    &info.URL,
		}
	})), nil
}

func (s *Server) GetProviders(ctx context.Context, _ api.GetProvidersRequestObject) (api.GetProvidersResponseObject, error) {
	providers := lo.Map(s.loaders, func(loader libmangal.ProviderLoader, _ int) api.Provider {
		info := loader.Info()
		return api.Provider{
			Description: &info.Description,
			Id:          info.ID,
			Name:        info.Name,
			Version:     info.Version,
		}
	})

	return api.GetProviders200JSONResponse(providers), nil
}

func NewServer() (*echo.Echo, error) {
	sub, err := fs.Sub(frontend, filepath.Join("ui", "dist"))
	if err != nil {
		return nil, err
	}

	server := &Server{}
	store, err := bigcache.NewStore(bigcache.Options{
		HardMaxCacheSize: 0,
		Eviction:         0,
		Codec:            encoding.Gob,
	})
	if err != nil {
		return nil, err
	}

	server.imageCache = store
	server.loaders, err = manager.Loaders()
	if err != nil {
		return nil, err
	}

	server.loadersByID = make(map[string]libmangal.ProviderLoader, len(server.loaders))
	for _, loader := range server.loaders {
		server.loadersByID[loader.Info().ID] = loader
	}

	handler := api.NewStrictHandler(server, nil)
	e := echo.New()
	api.RegisterHandlersWithBaseURL(e, handler, "api")

	e.StaticFS("/", sub)
	e.HideBanner = true

	return e, nil
}
