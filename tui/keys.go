package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/metafates/mangal/config"
)

// keyMap is a map of key bindings for the bubble.
type keyMap struct {
	state bubbleState

	Quit      key.Binding
	ForceQuit key.Binding
	Select    key.Binding
	SelectAll key.Binding
	Confirm   key.Binding
	Open      key.Binding
	Read      key.Binding
	Retry     key.Binding
	Back      key.Binding
	Filter    key.Binding

	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding

	Top    key.Binding
	Bottom key.Binding

	Help key.Binding
}

// shortHelpFor returns a short list of key bindings for the given state.
func (k keyMap) shortHelpFor(state bubbleState) []key.Binding {
	switch state {
	case ResumeState:
		return []key.Binding{k.Select, k.Confirm, k.Back, k.Filter}
	case SearchState:
		return []key.Binding{k.Confirm, k.ForceQuit}
	case LoadingState:
		return []key.Binding{k.ForceQuit}
	case MangaState:
		return []key.Binding{k.Open, k.Select, k.Back}
	case ChaptersState:
		return []key.Binding{k.Open, k.Read, k.Select, k.SelectAll, k.Confirm, k.Back}
	case ConfirmState:
		return []key.Binding{k.Confirm, k.Back, k.Quit}
	case DownloadingState:
		return []key.Binding{k.ForceQuit}
	case ExitPromptState:
		k.Open.SetHelp("o", "open folder")
		k.Retry.SetHelp("r", "redownload failed")
		return []key.Binding{k.Back, k.Open, k.Retry, k.Quit}
	}

	return []key.Binding{k.ForceQuit}
}

// fulleHelpFor returns a full list of key bindings for the given state.
func (k keyMap) fullHelpFor(state bubbleState) []key.Binding {
	switch state {
	case ResumeState:
		return []key.Binding{k.Select, k.Confirm, k.Back, k.Filter}
	case SearchState:
		return []key.Binding{k.Confirm, k.ForceQuit}
	case LoadingState:
		return []key.Binding{k.ForceQuit}
	case MangaState:
		k.Open.SetHelp("o", "open manga url")
		return []key.Binding{k.Open, k.Select, k.Back}
	case ChaptersState:
		k.Read.SetHelp("r", fmt.Sprintf("read chapter in the default %s app", string(config.UserConfig.Formats.Default)))
		k.Open.SetHelp("o", "open chapter url")
		k.Confirm.SetHelp("enter", "download selected chapters")
		return []key.Binding{k.Open, k.Read, k.Select, k.SelectAll, k.Confirm, k.Back}
	case ConfirmState:
		return []key.Binding{k.Confirm, k.Back, k.Quit}
	case DownloadingState:
		return []key.Binding{k.ForceQuit}
	case ExitPromptState:
		k.Open.SetHelp("o", "open folder")
		k.Retry.SetHelp("r", "redownload failed")
		return []key.Binding{k.Back, k.Open, k.Retry, k.Quit}
	}

	return []key.Binding{k.ForceQuit}
}

// ShortHelp returns a short list of key bindings for the given state.
func (k keyMap) ShortHelp() []key.Binding {
	return k.shortHelpFor(k.state)
}

// FullHelp returns a full list of key bindings for the given state.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.fullHelpFor(k.state)}
}
