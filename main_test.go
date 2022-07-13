package main

import (
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/filesystem"
	"github.com/spf13/afero"
)

func init() {
	// set memory filesystem for testing
	filesystem.Set(afero.NewMemMapFs())
	config.Initialize("", false)
}
