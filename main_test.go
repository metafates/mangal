package main

import (
	"github.com/metafates/mangal/filesystem"
	"github.com/spf13/afero"
)

func init() {
	// set memory filesystem for testing
	filesystem.Set(afero.NewMemMapFs())

}
