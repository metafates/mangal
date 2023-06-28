package model

import (
	"context"
	"github.com/Inno-Gang/goodle-cli/tui/base"
	"github.com/Inno-Gang/goodle-cli/tui/state"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/pkg/errors"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resize(base.Size{
			Width:  msg.Width,
			Height: msg.Height,
		})

		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.Back) && m.state.Backable():
			return m, m.back()
		case key.Matches(msg, m.keyMap.Help):
			m.help.ShowAll = !m.help.ShowAll
			m.resize(m.size)
			return m, nil
		}
	case base.MsgBack:
		// this msg can override Backable() output
		return m, m.back()
	case base.State:
		return m, m.pushState(msg)
	case error:
		if errors.Is(msg, context.Canceled) {
			return m, nil
		}

		log.Error(msg)

		return m, m.pushState(state.NewError(msg))
	}

	cmd := m.state.Update(m, msg)
	return m, cmd
}
