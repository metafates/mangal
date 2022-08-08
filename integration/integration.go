package integration

import (
	"github.com/metafates/mangal/integration/anilist"
	"github.com/metafates/mangal/source"
)

type Integrator interface {
	MarkRead(chapter *source.Chapter) error
}

var (
	Anilist Integrator = anilist.New()
)
