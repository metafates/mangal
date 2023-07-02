package base

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

type State interface {
	Intermediate() bool
	Backable() bool

	KeyMap() help.KeyMap

	Title() Title
	Subtitle() string
	Status() string

	Resize(size Size)

	Update(model Model, msg tea.Msg) tea.Cmd
	View(model Model) string
	Init(model Model) tea.Cmd
}
