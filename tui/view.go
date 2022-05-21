package tui

import (
	"fmt"
)

func plural(count int, word string) string {
	if count == 1 {
		return word
	}

	return word + "s"
}

// View renders the UI.
func (b Bubble) View() string {
	var rendered string

	switch b.state {
	case searchState:
		rendered = titleStyle.Render("Mangai")
		rendered += "\n\n" + b.input.View()
		rendered += "\n\n" + b.help.View(b.keys[b.state])
	case spinnerState:
		rendered = b.spinner.View() + " Searching..."
		rendered += "\n\n" + b.help.View(b.keys[b.state])
	case mangaSelectState:
		rendered = b.manga.View()
	case chaptersSelectState:
		rendered = b.chapters.View()
	case promptState:
		// TODO: make it better
		count := len(b.selected)
		rendered = fmt.Sprintf("Download %d %s of %s?", count, plural(count, "chapter"), b.prevManga)
		rendered += "\n\n" + b.help.View(b.keys[b.state])
	case progressState:
		// TODO: make it better
		rendered = titleStyle.Render("Downloading "+b.prevManga) + "\n\n"
		rendered += fmt.Sprintf(`%s

%s


%s

%s`, b.progress.View(), textSecondaryStyle.Render(b.prevChapter), b.subProgress.View(), textSecondaryStyle.Render(b.prevPanel))
		//if b.converting {
		//	rendered += fmt.Sprintf("%s\n%s\n\n%s\n%s", b.progress.View(), marginYStyle.Render(b.prevChapter), b.subProgress.View(), "Converting to pdf...")
		//} else {
		//	rendered += fmt.Sprintf("%s\n%s\n\n%s\n%s", b.progress.View(), b.prevChapter, b.subProgress.View(), b.prevPanel)
		//}
		rendered += "\n\n" + b.help.View(b.keys[b.state])
	case exitPrompt:
		count := len(b.selected)
		rendered = titleStyle.Render(fmt.Sprintf("%d %s of %s downloaded", count, plural(count, "chapter"), b.prevManga))
		rendered += "\n\n" + b.help.View(b.keys[b.state])
	}

	return commonStyle.Render(rendered)
}
