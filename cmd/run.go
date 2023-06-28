package cmd

import (
	"github.com/mangalorg/mangal/tui"
	"github.com/mangalorg/mangal/tui/state/loading"
)

type runCmd struct {
}

func (r *runCmd) Run() error {
	return tui.Run(loading.New("Hi mom"))
}
