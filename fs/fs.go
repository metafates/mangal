package fs

import "github.com/spf13/afero"

var FS = afero.Afero{
    Fs: afero.NewOsFs(),
}
