package integration

import (
	"github.com/metafates/mangal/integration/anilistintegration"
	"github.com/metafates/mangal/source"
)

type Integrator interface {
	MarkRead(chapter *source.Chapter) error
}

var (
	Anilist Integrator = anilistintegration.New()
)
