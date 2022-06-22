package main

import "github.com/spf13/afero"

// Afero is afero filesystem abstraction layer
var Afero = afero.Afero{Fs: afero.NewOsFs()}
