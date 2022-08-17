package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/provider"
	"github.com/spf13/viper"
)

func (b *statefulBubble) Init() tea.Cmd {
	if name := viper.GetString(config.DownloaderDefaultSource); name != "" {
		p, ok := provider.Get(name)
		if !ok {
			b.lastError = fmt.Errorf("provider %s not found", name)
			b.plot = randomPlot()
			b.newState(errorState)
			return nil
		}

		b.setState(loadingState)
		return tea.Batch(b.startLoading(), b.loadSource(p), b.waitForSourceLoaded())
	}

	return tea.Batch(textinput.Blink, b.loadSources())
}
