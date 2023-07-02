package main

import (
	"github.com/charmbracelet/log"
	"github.com/mangalorg/mangal/cmd"
	"github.com/mangalorg/mangal/config"
	"io"
)

func main() {
	if err := config.Load(); err != nil {
		log.Fatal(err)
	}

	// TODO: change this
	log.Default().SetOutput(io.Discard)
	cmd.Run()
}
