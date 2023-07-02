package chapsdownloaded

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mangalorg/libmangal"
)

type Options struct {
	Succeed, Failed  []libmangal.Chapter
	SucceedPaths     []string
	DownloadChapters func(chapters []libmangal.Chapter) tea.Cmd
}
