package tui

import "fmt"

func (b *statefulBubble) View() string {
	switch b.state {
	case idle:
		return b.viewIdle()
	case loadingState:
		return b.viewLoading()
	case historyState:
		return b.viewHistory()
	case sourcesState:
		return b.viewSources()
	case searchState:
		return b.viewSearch()
	case mangasState:
		return b.viewMangas()
	case chaptersState:
		return b.viewChapters()
	case confirmState:
		return b.viewConfirm()
	case readState:
		return b.viewRead()
	case downloadState:
		return b.viewDownload()
	case downloadDoneState:
		return b.viewDownloadDone()
	case exitState:
		return b.viewExit()
	case errorState:
		return b.viewError()
	}

	panic("unknown state")
}

func (b *statefulBubble) viewIdle() string {
	return ""
}

func (b *statefulBubble) viewLoading() string {
	return fmt.Sprintf("%s Loading...", b.spinnerC.View())
}

func (b *statefulBubble) viewHistory() string {
	return b.historyC.View()
}

func (b *statefulBubble) viewSources() string {
	return b.sourcesC.View()
}

func (b *statefulBubble) viewSearch() string {
	return b.inputC.View() + "\n" + b.helpC.View(b.keymap)
}

func (b *statefulBubble) viewMangas() string {
	return b.mangasC.View()
}

func (b *statefulBubble) viewChapters() string {
	return b.chaptersC.View()
}

func (b *statefulBubble) viewConfirm() string {
	return fmt.Sprintf("Download %d chapters?\n%s", len(b.selectedChapters), b.helpC.View(b.keymap))
}

func (b *statefulBubble) viewRead() string {
	return b.spinnerC.View() + b.progressStatus + "\n" + b.helpC.View(b.keymap)
}

func (b *statefulBubble) viewDownload() string {
	return b.progressC.View() + "\n" + b.spinnerC.View() + b.progressStatus + "\n" + b.helpC.View(b.keymap)
}

func (b *statefulBubble) viewDownloadDone() string {
	return "Download done" + "\n" + b.helpC.View(b.keymap)
}

func (b *statefulBubble) viewExit() string {
	return ""
}

func (b *statefulBubble) viewError() string {
	return "Error"
}
