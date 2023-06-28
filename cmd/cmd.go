package cmd

import "github.com/alecthomas/kong"

var cli struct {
    Version versionCmd `cmd:"" help:"Print version"`    
}

func Run() {
    ctx := kong.Parse(&cli)
    err := ctx.Run()
    ctx.FatalIfErrorf(err)
}
