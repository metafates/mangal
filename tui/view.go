package tui

import (
	"fmt"
	"github.com/metafates/mangal/icon"
	"github.com/metafates/mangal/style"
	"github.com/metafates/mangal/util"
	"math/rand"
	"strings"
)

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
	case errorState:
		return b.viewError()
	}

	panic("unknown state")
}

func (b *statefulBubble) viewIdle() string {
	return ""
}

func (b *statefulBubble) viewLoading() string {
	return b.renderLines(true, b.spinnerC.View()+"Loading...")
}

func (b *statefulBubble) viewHistory() string {
	return b.historyC.View()
}

func (b *statefulBubble) viewSources() string {
	return b.sourcesC.View()
}

func (b *statefulBubble) viewSearch() string {
	return b.renderLines(true, b.inputC.View())
}

func (b *statefulBubble) viewMangas() string {
	return b.mangasC.View()
}

func (b *statefulBubble) viewChapters() string {
	return b.chaptersC.View()
}

func (b *statefulBubble) viewConfirm() string {
	return b.renderLines(
		true,
		fmt.Sprintf("Download %d chapters?", len(b.selectedChapters)),
	)
}

func (b *statefulBubble) viewRead() string {
	return b.renderLines(
		true,
		b.progressC.View(),
		b.spinnerC.View()+b.progressStatus,
	)
}

func (b *statefulBubble) viewDownload() string {
	return b.renderLines(
		true,
		b.progressC.View(),
		b.spinnerC.View()+b.progressStatus,
	)
}

func (b *statefulBubble) viewDownloadDone() string {
	return b.renderLines(true, "Download finished")
}

func (b *statefulBubble) viewError() string {
	return b.renderLines(
		true,
		icon.Get(icon.Fail)+" Uggh, something went wrong. Maybe try again?",
		"",
		style.Italic(util.Wrap(randomPlot(), b.terminalWidth/2)),
		"",
		style.Combined(style.Italic, style.Red)(b.lastError.Error()),
	)
}

func (b *statefulBubble) renderLines(addHelp bool, lines ...string) string {
	l := strings.Join(lines, "\n")
	if addHelp {
		l += "\n\n" + b.helpC.View(b.keymap)
	}

	return l
}

func randomPlot() string {
	plots := []string{
		"The universe is a dangerous place. There are many things that can go wrong. This is one of them:",
		"Heroically fighting an endless army of errors and bugs Mangal died a hero. Their last words were:",
	}

	return plots[rand.Intn(len(plots))]
}
