package manager

import (
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/path"
	"github.com/mangalorg/mangal/provider/bundle"
)

func Loaders() ([]libmangal.ProviderLoader, error) {
	return bundle.Loaders(path.ProvidersDir())
}
