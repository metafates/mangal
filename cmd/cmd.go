package cmd

import (
	"flag"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/metafates/mangai/tui"
)

func Execute(version string, build string) error {
	var showVersion bool

	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.Parse()

	if showVersion {
		fmt.Printf("Mangai version %s build %s\n", version, build)
		return nil
	}

	program := tea.NewProgram(tui.New(), tea.WithAltScreen())
	return program.Start()
}
