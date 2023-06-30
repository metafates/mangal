package cmd

import (
	"github.com/alecthomas/kong"
)

type rootCmd struct {
	Run     runCmd     `cmd:"" help:"Run mangal" default:"1"`
	Version versionCmd `cmd:"" help:"Print version"`
	Path    pathCmd    `cmd:"" help:"Print paths"`
}

var cli rootCmd

func Run() {
	ctx := kong.Parse(&cli, kong.ShortUsageOnError())
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
