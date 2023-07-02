package util

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/lipgloss"
)

func NewList[T any](
	delegateHeight int,
	singular, plural string,
	items []T,
	transform func(T) list.DefaultItem,
) list.Model {
	var listItems = make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = transform(item)
	}

	border := lipgloss.ThickBorder()

	delegate := list.NewDefaultDelegate()

	delegate.Styles.NormalTitle.Bold(true)
	delegate.Styles.SelectedTitle.Bold(true)
	delegate.Styles.SelectedTitle.Border(border, false, false, false, true)
	delegate.Styles.SelectedDesc.
		Border(border, false, false, false, true).
		Foreground(delegate.Styles.NormalDesc.GetForeground())

	//delegate.Styles.SelectedTitle.BorderLeftForeground(color.Accent)
	//delegate.Styles.SelectedDesc.BorderLeftForeground(color.Accent)

	if delegateHeight == 1 {
		delegate.ShowDescription = false
	}

	delegate.SetHeight(delegateHeight)

	l := list.New(listItems, delegate, 0, 0)
	l.SetShowHelp(false)
	l.SetShowFilter(false)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	l.SetShowPagination(false)
	l.InfiniteScrolling = true
	l.KeyMap.CancelWhileFiltering = Bind("cancel", "esc")

	l.Paginator.Type = paginator.Arabic

	l.SetStatusBarItemName(singular, plural)

	return l
}
