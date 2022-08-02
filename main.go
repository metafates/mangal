package main

import (
	"github.com/metafates/mangal/cmd"
	"github.com/metafates/mangal/config"
	"github.com/samber/lo"
)

func main() {
	lo.Must0(config.Setup())
	cmd.Execute()
}
