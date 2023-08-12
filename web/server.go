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
	imageCache gokv.Store
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
	loaders, err := manager.Loaders()
	if err != nil {
		return nil, err
	}

	loader, ok := lo.Find(loaders, func(loader libmangal.ProviderLoader) bool {
		return loader.Info().ID == request.Params.Provider
	})

	if !ok {
		return api.SearchMangasdefaultJSONResponse{
			StatusCode: 404,
			Body: api.Error{
				Code:    404,
				Message: fmt.Sprintf("Provider %q not found", request.Params.Provider),
			},
		}, nil
	}

	clientInstance, err := client.NewClient(ctx, loader)
	if err != nil {
		return api.SearchMangasdefaultJSONResponse{
			StatusCode: 400,
			Body: api.Error{
				Message: err.Error(),
			},
		}, nil
	}
	defer clientInstance.Close()

	mangas, err := clientInstance.SearchMangas(ctx, request.Params.Query)
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
	loaders, err := manager.Loaders()
	if err != nil {
		return nil, err
	}

	providers := lo.Map(loaders, func(loader libmangal.ProviderLoader, _ int) api.Provider {
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

	handler := api.NewStrictHandler(server, nil)
	e := echo.New()
	api.RegisterHandlersWithBaseURL(e, handler, "api")

	e.StaticFS("/", sub)
	e.HideBanner = true

	return e, nil
}
