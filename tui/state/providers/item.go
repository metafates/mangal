package providers

import (
	"fmt"

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
	return fmt.Sprintf(`Version %s
%s`,
		i.Info().Version,
		i.Info().Website,
	)
}
