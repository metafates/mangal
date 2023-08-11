package web

import (
	"context"
	"embed"
	"io/fs"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/meta"
	"github.com/mangalorg/mangal/provider/manager"
	"github.com/mangalorg/mangal/web/api"
	"github.com/samber/lo"
)

//go:embed all:ui/dist
var frontend embed.FS

var _ api.StrictServerInterface = (*Server)(nil)

type Server struct {
	loaders  []libmangal.ProviderLoader
	provider libmangal.Provider
}

// GetMangalInfo implements api.StrictServerInterface.
func (*Server) GetMangalInfo(ctx context.Context, request api.GetMangalInfoRequestObject) (api.GetMangalInfoResponseObject, error) {
	return api.GetMangalInfo200JSONResponse{
		Version: meta.Version,
	}, nil
}

// SearchMangas implements api.StrictServerInterface.
func (s *Server) SearchMangas(ctx context.Context, request api.SearchMangasRequestObject) (api.SearchMangasResponseObject, error) {
	loader, ok := lo.Find(s.loaders, func(loader libmangal.ProviderLoader) bool {
		return loader.Info().ID == request.Params.ProviderId
	})

	if !ok {
		return api.SearchMangasdefaultJSONResponse{
			StatusCode: 404,
			Body: api.Error{
				Code:    404,
				Message: "Provider not found",
			},
		}, nil
	}

	provider, err := loader.Load(ctx)
	if err != nil {
		return api.SearchMangasdefaultJSONResponse{
			StatusCode: 400,
			Body: api.Error{
				Message: err.Error(),
			},
		}, nil
	}

	s.provider = provider

	mangas, err := provider.SearchMangas(ctx, request.Params.Query)
	if err != nil {
		return api.SearchMangasdefaultJSONResponse{
			StatusCode: 400,
			Body: api.Error{
				Message: err.Error(),
			},
		}, nil
	}

	return api.SearchMangas200JSONResponse(lo.Map(mangas, func(manga libmangal.Manga, _ int) api.Manga {
		return api.Manga{
			Title: manga.Info().Title,
		}
	})), nil
}

func (s *Server) GetProviders(ctx context.Context, _ api.GetProvidersRequestObject) (api.GetProvidersResponseObject, error) {
	loaders, err := manager.Loaders()
	if err != nil {
		return nil, err
	}

	s.loaders = loaders

	providers := lo.Map(loaders, func(loader libmangal.ProviderLoader, _ int) api.Provider {
		info := loader.Info()
		return api.Provider{
			Id:   info.ID,
			Name: &info.Name,
		}
	})

	return api.GetProviders200JSONResponse(providers), nil
}

func NewServer() (*echo.Echo, error) {
	sub, err := fs.Sub(frontend, filepath.Join("ui", "dist"))
	if err != nil {
		return nil, err
	}

	handler := api.NewStrictHandler(&Server{}, nil)
	e := echo.New()
	api.RegisterHandlersWithBaseURL(e, handler, "api")

	e.StaticFS("/", sub)

	return e, nil
}
