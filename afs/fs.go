package afs

import "github.com/spf13/afero"

var Afero = afero.Afero{
	Fs: afero.NewOsFs(),
}
