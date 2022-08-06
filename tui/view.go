package tui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
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
	return b.renderLines(
		true,
		[]string{
			b.spinnerC.View() + "Loading...",
		},
	)
}

func (b *statefulBubble) viewHistory() string {
	return b.historyC.View()
}

func (b *statefulBubble) viewSources() string {
	return b.sourcesC.View()
}

func (b *statefulBubble) viewSearch() string {
	return b.renderLines(
		true,
		[]string{b.inputC.View()},
	)
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
		[]string{
			fmt.Sprintf(icon.Get(icon.Question)+" Download %d chapters?", len(b.selectedChapters)),
		},
	)
}

func (b *statefulBubble) viewRead() string {
	var chapterName string

	chapter := b.currentDownloadingChapter
	if chapter != nil {
		chapterName = chapter.Name
	}

	return b.renderLines(
		true,
		[]string{
			style.Trim(b.width)(fmt.Sprintf(icon.Get(icon.Progress)+" Downloading chapter %s", style.Magenta(chapterName))),
			"",
			style.Trim(b.width)(b.spinnerC.View() + b.progressStatus),
		},
	)
}

func (b *statefulBubble) viewDownload() string {
	var chapterName string

	chapter := b.currentDownloadingChapter
	if chapter != nil {
		chapterName = chapter.Name
	}

	return b.renderLines(
		true,
		[]string{
			style.Trim(b.width)(fmt.Sprintf(icon.Get(icon.Progress)+" Downloading chapter %s", style.Magenta(chapterName))),
			"",
			b.progressC.View(),
			"",
			style.Trim(b.width)(b.spinnerC.View() + b.progressStatus),
		},
	)
}

func (b *statefulBubble) viewDownloadDone() string {
	return b.renderLines(
		true,
		[]string{icon.Get(icon.Success) + " Download finished. *Beep-Boop-Boop*"},
	)
}

func (b *statefulBubble) viewError() string {
	errorMsg := util.Wrap(style.Combined(style.Italic, style.Red)(b.lastError.Error()), b.width)
	return b.renderLines(
		true,
		append([]string{
			icon.Get(icon.Fail) + " Uggh, something went wrong. Maybe try again?",
			"",
			style.Italic(util.Wrap(b.plot, b.width)),
			"",
		},
			strings.Split(errorMsg, "\n")...,
		),
	)
}

var (
	paddingStyle = lipgloss.NewStyle().PaddingTop(1).PaddingLeft(2).PaddingRight(2)
)

func (b *statefulBubble) renderLines(addHelp bool, lines []string) string {
	h := len(lines)
	l := strings.Join(lines, "\n")
	if addHelp {
		l += strings.Repeat("\n", b.height-h) + b.helpC.View(b.keymap)
	}

	return paddingStyle.Render(l)
}

func randomPlot() string {
	plots := []string{
		"The universe is a dangerous place. There are many things that can go wrong. This is one of them:",
		"Fighting an endless army of errors and bugs Mangal died a hero. Their last words were:",
		"I used to download stuff without any errors, then I took an arrow to the knee. By arrow I mean this:",
	}

	return plots[rand.Intn(len(plots))]
}
