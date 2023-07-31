package providers

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
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
	return fmt.Sprint(i.FilterValue(), " ", lipgloss.NewStyle().Italic(true).Render("v"+i.Info().Version))
}

func (i Item) Description() string {
	return i.Info().Website
}
