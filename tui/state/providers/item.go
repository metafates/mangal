package providers

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/mangalorg/libmangal"
)

var _ list.Item = (*Item)(nil)

type Item struct {
	libmangal.ProviderLoader
}

func (i Item) FilterValue() string {
	return i.String()
}

func (i Item) Title() string {
	return i.FilterValue()
}

func (i Item) Description() string {
	return i.Info().Website
}
