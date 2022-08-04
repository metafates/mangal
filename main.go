package main

import (
	"github.com/metafates/mangal/cmd"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/icon"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

func main() {
	lo.Must0(config.Setup())
	icon.SetVariant(viper.GetString(config.IconsVariant))
	cmd.Execute()
}
