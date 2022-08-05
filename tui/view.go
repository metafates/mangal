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
	case readDownloadState:
		return b.viewReadDownload()
	case readDownloadDoneState:
		return b.viewReadDownloadDone()
	case downloadState:
		return b.viewDownload()
	case downloadDoneState:
		return b.viewDownloadDone()
	case exitState:
		return b.viewExit()
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
	return b.inputC.View()
}

func (b *statefulBubble) viewMangas() string {
	return b.mangasC.View()
}

func (b *statefulBubble) viewChapters() string {
	return b.chaptersC.View()
}

func (b *statefulBubble) viewConfirm() string {
	return ""
}

func (b *statefulBubble) viewReadDownload() string {
	return ""
}

func (b *statefulBubble) viewReadDownloadDone() string {
	return ""
}

func (b *statefulBubble) viewDownload() string {
	return ""
}

func (b *statefulBubble) viewDownloadDone() string {
	return ""
}

func (b *statefulBubble) viewExit() string {
	return ""
}
