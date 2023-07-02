package chapters

import (
	"fmt"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/icon"
	"github.com/zyedidia/generic/set"
)

type Item struct {
	chapter       libmangal.Chapter
	selectedItems *set.Set[*Item]
}

func (i *Item) FilterValue() string {
	return i.chapter.String()
}

func (i *Item) Title() string {
	if i.IsSelected() {
		return fmt.Sprint(i.FilterValue(), " ", icon.Mark)
	}

	return i.FilterValue()
}

func (i *Item) Description() string {
	return i.chapter.Info().URL
}

func (i *Item) IsSelected() bool {
	return i.selectedItems.Has(i)
}

func (i *Item) Toggle() {
	if i.IsSelected() {
		i.selectedItems.Remove(i)
	} else {
		i.selectedItems.Put(i)
	}
}
