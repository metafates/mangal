package main

import (
	"github.com/mangalorg/mangal/cmd"
	"github.com/mangalorg/mangal/config"
	"github.com/mangalorg/mangal/log"
)

func main() {
	if err := config.Load(); err != nil {
		log.L.Fatal().Err(err).Msg("failed to load config")
		panic(err)
	}

	cmd.Execute()
}
