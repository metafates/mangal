package chapters

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/tui/state/loading"
)

func downloadChapterCmd(
	ctx context.Context,
	client *libmangal.Client,
	chapter libmangal.Chapter,
	options libmangal.DownloadOptions,
	onSuccess func(path string) tea.Msg,
) tea.Cmd {
	loadingState := loading.New("Preparing...")
	return tea.Sequence(
		func() tea.Msg {
			return loadingState
		},
		func() tea.Msg {
			client.SetLogFunc(func(msg string) {
				loadingState.SetMessage(msg)
			})

			chapterPath, err := client.DownloadChapter(ctx, chapter, options)
			if err != nil {
				return err
			}

			return onSuccess(chapterPath)
		},
	)
}
