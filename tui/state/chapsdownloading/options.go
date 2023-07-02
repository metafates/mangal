package chapsdownloading

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mangalorg/libmangal"
)

type Options struct {
	DownloadChapter    func(ctx context.Context, chapter libmangal.Chapter) (string, error)
	OnDownloadFinished func(paths []string, succeed, failed []libmangal.Chapter) tea.Cmd
}
