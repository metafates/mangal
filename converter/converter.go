package converter

import (
	"fmt"
	"github.com/metafates/mangal/converter/cbz"
	"github.com/metafates/mangal/converter/pdf"
	"github.com/metafates/mangal/converter/plain"
	"github.com/metafates/mangal/converter/zip"
	"github.com/metafates/mangal/source"
)

type Converter interface {
	Save(chapter *source.Chapter) (string, error)
	SaveTemp(chapter *source.Chapter) (string, error)
}

const (
	Plain = "plain"
	CBZ   = "cbz"
	PDF   = "pdf"
	ZIP   = "zip"
)

var converters = map[string]Converter{
	Plain: plain.New(),
	CBZ:   cbz.New(),
	PDF:   pdf.New(),
	ZIP:   zip.New(),
}

func Available() []string {
	return []string{
		Plain,
		CBZ,
		PDF,
		ZIP,
	}
}

func Get(name string) (Converter, error) {
	if converter, ok := converters[name]; ok {
		return converter, nil
	}

	return nil, fmt.Errorf("unkown format \"%s\"", name)
}
