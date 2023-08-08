package chapsdownloaded

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mangalorg/mangal/color"
	"github.com/mangalorg/mangal/stringutil"
	"github.com/mangalorg/mangal/tui/base"
	"github.com/skratchdot/open-golang/open"
	"path/filepath"
	"strings"
)

var _ base.State = (*State)(nil)

type State struct {
	options Options
	keyMap  KeyMap
}

func (s *State) Intermediate() bool {
	return true
}

func (s *State) KeyMap() help.KeyMap {
	return s.keyMap
}

func (s *State) Title() base.Title {
	return base.Title{Text: "Done"}
}

func (s *State) Subtitle() string {
	return ""
}

func (s *State) Status() string {
	return ""
}

func (s *State) Backable() bool {
	return true
}

func (s *State) Resize(size base.Size) {
}

func (s *State) Update(model base.Model, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keyMap.Quit):
			return tea.Quit
		case key.Matches(msg, s.keyMap.Open) && len(s.options.SucceedPaths) > 0:
			err := open.Start(filepath.Dir(s.options.SucceedPaths[0]))
			if err != nil {
				return func() tea.Msg {
					return err
				}
			}

			return nil
		case key.Matches(msg, s.keyMap.Retry) && len(s.options.Failed) > 0:
			return s.options.DownloadChapters(s.options.Failed)
		}
	}

	return nil
}

func (s *State) View(model base.Model) string {
	var (
		succeed = len(s.options.Succeed)
		failed  = len(s.options.Failed)
	)

	if len(s.options.Failed) == 0 {
		return lipgloss.
			NewStyle().
			Foreground(color.Success).
			Render(fmt.Sprintf(
				"%s downloaded successfully!",
				stringutil.Quantify(succeed, "chapter", "chapters"),
			))
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(
		"%s downloaded successfully, %d failed.",
		stringutil.Quantify(succeed, "chapter", "chapters"),
		failed,
	))

	sb.WriteString("\n\nFailed:\n")

	if failed <= 3 {
		for _, chapter := range s.options.Failed {
			sb.WriteString(fmt.Sprintf("\n%s", chapter))
		}
	} else {
		var indices = make([]float32, failed)
		for i, c := range s.options.Failed {
			indices[i] = c.Info().Number
		}

		sb.WriteString(stringutil.FormatRanges(indices))
	}

	return sb.String()
}

func (s *State) Init(model base.Model) tea.Cmd {
	return nil
}
