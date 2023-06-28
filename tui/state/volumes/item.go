package volumes

import (
	"fmt"
	"github.com/mangalorg/libmangal"
)

type Item struct {
	libmangal.Volume
}

func (i Item) FilterValue() string {
	return fmt.Sprintf("Volume %d", i.Info().Number)
}

func (i Item) Title() string {
	return i.FilterValue()
}

func (i Item) Description() string {
	return ""
}
