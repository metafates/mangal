package main

import (
	"github.com/charmbracelet/log"
	"github.com/mangalorg/mangal/cmd"
	"io"
)

func main() {
	// TODO: change this
	log.Default().SetOutput(io.Discard)
	cmd.Run()
}
