package main

import (
	"io"

	"github.com/charmbracelet/log"
	"github.com/mangalorg/mangal/cmd"
	"github.com/mangalorg/mangal/config"
)

func main() {
	if err := config.Load(); err != nil {
		log.Fatal("failed to load config", "err", err)
	}

	// TODO: change this
	log.Default().SetOutput(io.Discard)
	cmd.Execute()
}
