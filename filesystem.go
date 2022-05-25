package main

import "github.com/spf13/afero"

var FileSystem = afero.NewOsFs()
var Afero = afero.Afero{Fs: FileSystem}
