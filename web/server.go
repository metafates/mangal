package web

import (
	"context"
	"embed"
	"io/fs"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/provider/manager"
	"github.com/mangalorg/mangal/web/api"
	"github.com/samber/lo"
)

//go:embed all:ui/build
var frontend embed.FS

var _ api.StrictServerInterface = (*Server)(nil)

type Server struct{}

func (s *Server) GetProviders(ctx context.Context, _ api.GetProvidersRequestObject) (api.GetProvidersResponseObject, error) {
	loaders, err := manager.Loaders()
	if err != nil {
		return nil, err
	}

	return api.GetProviders200JSONResponse(lo.Map(loaders, func(l libmangal.ProviderLoader, _ int) api.Provider {
		info := l.Info()
		return api.Provider{
			Id:   info.ID,
			Name: &info.Name,
		}
	})), nil
}

func NewServer() (*echo.Echo, error) {

	sub, err := fs.Sub(frontend, filepath.Join("ui", "build"))
	if err != nil {
		return nil, err
	}

	handler := api.NewStrictHandler(&Server{}, nil)
	e := echo.New()
	api.RegisterHandlersWithBaseURL(e, handler, "api")

	e.StaticFS("/", sub)

	return e, nil
}
