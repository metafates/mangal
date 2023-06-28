package cmd

import (
	"github.com/mangalorg/mangal/provider"
	"github.com/mangalorg/mangal/tui"
	"github.com/mangalorg/mangal/tui/state/providers"
)

type runCmd struct {
}

func (r *runCmd) Run() error {
	loaders, err := provider.InstalledProviders()
	if err != nil {
		return err
	}

	return tui.Run(providers.New(loaders))
}
