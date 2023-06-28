package chapters

import (
	"github.com/mangalorg/libmangal"
)

type Item struct {
	libmangal.Chapter
}

func (i Item) FilterValue() string {
	return i.String()
}

func (i Item) Title() string {
	return i.FilterValue()
}

func (i Item) Description() string {
	return i.Info().URL
}
