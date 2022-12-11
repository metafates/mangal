package tui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/metafates/mangal/color"
	"github.com/metafates/mangal/icon"
	"github.com/metafates/mangal/key"
	"github.com/metafates/mangal/style"
	"github.com/metafates/mangal/util"
	"github.com/muesli/reflow/wrap"
	"github.com/spf13/viper"
	"math/rand"
	"strconv"
	"strings"
)

func (b *statefulBubble) View() string {
	switch b.state {
	case scrapersInstallState:
		return b.viewScrapersInstallState()
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
	case anilistSelectState:
		return b.viewAniList()
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

func (b *statefulBubble) viewLoading() string {
	return b.renderLines(
		true,
		[]string{
			style.Title("Loading"),
			"",
			b.spinnerC.View() + b.progressStatus,
		},
	)
}

func (b *statefulBubble) viewHistory() string {
	return listExtraPaddingStyle.Render(b.historyC.View())
}

func (b *statefulBubble) viewSources() string {
	return listExtraPaddingStyle.Render(b.sourcesC.View())
}

func (b *statefulBubble) viewSearch() string {
	lines := []string{
		style.Title("Search Manga"),
		"",
		b.inputC.View(),
	}

	if b.searchSuggestion.IsPresent() {
		lines = append(
			lines,
			"",
			fmt.Sprintf("Search %s ?", style.Fg(color.Orange)(b.searchSuggestion.MustGet())),
			"",
			fmt.Sprintf("Press %s to accept", style.Bold(style.Faint(b.keymap.acceptSearchSuggestion.Help().Key))),
		)
	}

	return b.renderLines(
		true,
		lines,
	)
}

func (b *statefulBubble) viewMangas() string {
	return listExtraPaddingStyle.Render(b.mangasC.View())
}

func (b *statefulBubble) viewChapters() string {
	return listExtraPaddingStyle.Render(b.chaptersC.View())
}

func (b *statefulBubble) viewAniList() string {
	return listExtraPaddingStyle.Render(b.anilistC.View())
}

func (b *statefulBubble) viewConfirm() string {
	return b.renderLines(
		true,
		[]string{
			style.Title("Confirm"),
			"",
			fmt.Sprintf("%s Download %s?", icon.Get(icon.Question), util.Quantify(len(b.selectedChapters), "chapter", "chapters")),
		},
	)
}

func (b *statefulBubble) downloadingChapterMetainfo() string {
	metainfo := strings.Builder{}

	// Even though when this function is called chapter isn't supposed to be nil,
	// it can be one for a brief moment.
	// I assume that it's because View() is called before Update()
	if b.currentDownloadingChapter != nil {
		metainfo.WriteString("From ")
		metainfo.WriteString(style.Fg(color.Orange)(b.currentDownloadingChapter.Source().Name()))
		metainfo.WriteString(" as ")
	}

	metainfo.WriteString(style.Fg(color.Purple)(viper.GetString(key.FormatsUse)))
	return metainfo.String()
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
			style.Title("Reading"),
			"",
			style.Truncate(b.width)(fmt.Sprintf(icon.Get(icon.Progress)+" Downloading %s", style.Fg(color.Purple)(chapterName))),
			"",
			style.Truncate(b.width)(b.spinnerC.View() + b.progressStatus),
			"",
			style.Truncate(b.width)(b.downloadingChapterMetainfo()),
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
			style.Title("Downloading"),
			"",
			style.Truncate(b.width)(fmt.Sprintf(icon.Get(icon.Progress)+" Downloading %s", style.Fg(color.Purple)(chapterName))),
			"",
			b.progressC.View(),
			"",
			style.Truncate(b.width)(b.spinnerC.View() + b.progressStatus),
			"",
			style.Truncate(b.width)(b.downloadingChapterMetainfo()),
		},
	)
}

func (b *statefulBubble) viewDownloadDone() string {
	failed := len(b.failedChapters)
	succeded := len(b.succededChapters)

	var msg string

	{
		temp := strings.Split(util.Quantify(succeded, "chapter", "chapters"), " ")
		temp[0] = style.Fg(color.Green)(temp[0])
		s := strings.Join(temp, " ") + " downloaded"
		f := fmt.Sprintf("%s failed", style.Fg(color.Red)(strconv.Itoa(failed)))

		msg = fmt.Sprintf("%s, %s", s, f)
	}

	lines := []string{
		style.Title("Finish"),
		"",
		msg,
	}

	if succeded > 0 && viper.GetBool(key.TUIShowDownloadedPath) {
		path, err := b.selectedManga.Path(false)
		if err == nil {
			lines = append(lines, "")
			lines = append(lines, fmt.Sprintf("Downloaded to %s", style.Faint(path)))
		}
	}

	return b.renderLines(
		true,
		lines,
	)
}

func (b *statefulBubble) viewError() string {
	errorMsg := wrap.String(style.New().Italic(true).Foreground(color.Red).Render(b.lastError.Error()), b.width)
	return b.renderLines(
		true,
		append([]string{
			style.ErrorTitle("Error"),
			"",
			icon.Get(icon.Fail) + " Uggh, something went wrong. Maybe try again?",
			"",
		},
			strings.Split(wrap.String(style.Italic(b.errorPlot), b.width)+"\n\n"+errorMsg, "\n")...,
		),
	)
}

func (b *statefulBubble) viewScrapersInstallState() string {
	return listExtraPaddingStyle.Render(b.scrapersInstallC.View())
}

var (
	listExtraPaddingStyle = lipgloss.NewStyle().Padding(1, 2, 1, 0)
	paddingStyle          = lipgloss.NewStyle().Padding(1, 2)
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
