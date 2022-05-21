package cmd

import (
	"flag"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/metafates/mangai/config"
	"github.com/metafates/mangai/shared"
	"github.com/metafates/mangai/tui"
)

func Execute(version string, build string) error {
	var showVersion bool

	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.Parse()

	if showVersion {
		fmt.Printf("%s version %s\nBuild %s\n", shared.Mangai, version, build)
		return nil
	}

	var program *tea.Program

	if config.Get().Fullscreen {
		program = tea.NewProgram(tui.New(), tea.WithAltScreen())
	} else {
		program = tea.NewProgram(tui.New())
	}

	return program.Start()
}
