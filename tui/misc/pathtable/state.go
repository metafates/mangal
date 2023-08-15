package pathtable

import (
	"fmt"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	table   table.Model
	keyMap  keyMap
	timer   *time.Timer
	help    help.Model
	msg     string
	quiting bool
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Quit):
			m.quiting = true
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.Copy):
			row := m.table.SelectedRow()
			path := row[1]

			if m.timer != nil {
				m.timer.Stop()
			}

			if err := clipboard.WriteAll(path); err != nil {
				m.msg = fmt.Sprintf("error: %s", err)
			} else {
				m.msg = fmt.Sprintf("copied %q path to clipboard", row[0])
			}

			delay := time.Second
			m.timer = time.AfterFunc(delay, func() {
				m.msg = ""
			})

			return m, func() tea.Msg {
				// force update
				time.Sleep(delay)
				return struct{}{}
			}
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	if m.quiting {
		return ""
	}

	style := lipgloss.NewStyle().Margin(1, 2)

	return style.Render(strings.Join([]string{
		m.table.View(),
		"",
		m.help.View(&m.keyMap),
		"",
		lipgloss.NewStyle().Italic(true).Faint(true).Render(m.msg),
	}, "\n"))
}
