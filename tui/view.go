package tui

import (
	"fmt"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/style"
	"github.com/metafates/mangal/util"
	"golang.org/x/exp/slices"
	"log"
	"strconv"
)

// View handles how the Bubble should be rendered
func (b *Bubble) View() string {
	var view string

	template := viewTemplates[b.state]

	switch b.state {
	case ResumeState:
		view = fmt.Sprintf(template, b.ResumeList.View())
	case SearchState:
		if config.UserConfig.UI.Title == "" {
			view = b.input.View()
		} else {
			view = fmt.Sprintf(template, style.InputTitleStyle.Render(config.UserConfig.UI.Title), b.input.View())
		}
	case LoadingState:
		view = fmt.Sprintf(template, b.spinner.View())
	case MangaState:
		view = fmt.Sprintf(template, b.mangaList.View())
	case ChaptersState:
		view = fmt.Sprintf(template, b.chaptersList.View())
	case ConfirmState:
		// Should be unreachable
		if len(b.selectedChapters) == 0 {
			log.Fatal("No chapters selected")
		}

		mangaName := b.chaptersList.Items()[0].(*listItem).url.Relation.Info
		chaptersToDownload := len(b.selectedChapters)
		view = fmt.Sprintf(
			template,
			style.AccentStyle.Render(strconv.Itoa(chaptersToDownload)),
			util.Plural("chapter", chaptersToDownload),
			style.AccentStyle.Render(util.PrettyTrim(mangaName, 40)),
		)
	case DownloadingState:

		var header string

		// It shouldn't be nil at this stage but it panics TODO: FIX THIS
		if b.chaptersDownloadProgressInfo.Current != nil {
			mangaName := b.chaptersDownloadProgressInfo.Current.Info
			header = fmt.Sprintf("Downloading %s", util.PrettyTrim(style.AccentStyle.Render(mangaName), 40))
		} else {
			header = "Preparing for download..."
		}

		subheader := b.chapterDownloadProgressInfo.Message
		view = fmt.Sprintf("%s\n\n%s\n\n%s %s", header, b.progress.View(), b.spinner.View(), subheader)
	case ExitPromptState:
		succeeded := b.chaptersDownloadProgressInfo.Succeeded
		failed := b.chaptersDownloadProgressInfo.Failed

		succeededRendered := style.SuccessStyle.Render(strconv.Itoa(len(succeeded)))
		failedRendered := style.FailStyle.Render(strconv.Itoa(len(failed)))

		view = fmt.Sprintf(template, succeededRendered, util.Plural("chapter", len(succeeded)), failedRendered)

		// show failed chapters
		for _, chapter := range failed {
			view += fmt.Sprintf("\n\n%s %s", style.FailStyle.Render("Failed"), chapter.Info)
		}
	}

	// Do not add help Bubble at these states, since they already have one
	if slices.Contains([]bubbleState{MangaState, ChaptersState, ResumeState}, b.state) {
		return style.CommonStyle.Render(view)
	}

	// Add help view
	return style.CommonStyle.Render(fmt.Sprintf("%s\n\n%s", view, b.help.View(b.keyMap)))
}

// viewTemplates is a map of the templates for the different states
var viewTemplates = map[bubbleState]string{
	ResumeState:      "%s",
	SearchState:      "%s\n\n%s",
	LoadingState:     "%s Searching...",
	MangaState:       "%s",
	ChaptersState:    "%s",
	ConfirmState:     "Download %s %s of %s ?",
	DownloadingState: "%s\n\n%s\n\n%s %s",
	ExitPromptState:  "Done. %s %s downloaded, %s failed",
}
