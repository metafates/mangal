package converter

import (
	"github.com/metafates/mangal/converter/cbz"
	"github.com/metafates/mangal/converter/plain"
	"github.com/metafates/mangal/source"
)

type Converter interface {
	Save(chapter *source.Chapter) (string, error)
}

var converters = map[string]Converter{
	"plain": plain.New(),
	"cbz":   cbz.New(),
}

func Converters() map[string]Converter {
	return converters
}
