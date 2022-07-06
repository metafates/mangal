package tui

import (
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/style"
)

// resize the bubble
func (b *Bubble) resize(width int, height int) {
	// Set size to minimum for non-fullscreen runtime
	if !config.UserConfig.UI.Fullscreen {
		b.mangaList.SetSize(0, 0)
		b.chaptersList.SetSize(0, 0)
		b.ResumeList.SetSize(0, 0)
		return
	}

	x, y := style.CommonStyle.GetFrameSize()
	b.mangaList.SetSize(width-x, height-y)
	b.chaptersList.SetSize(width-x, height-y)
	b.ResumeList.SetSize(width-x, height-y)
}

// setState sets the state of the bubble
func (b *Bubble) setState(state bubbleState) {
	b.state = state
	b.keyMap.state = state
}
