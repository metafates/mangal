package shared

import "github.com/spf13/afero"

var AferoBackend = afero.NewOsFs()
var AferoFS = afero.Afero{Fs: AferoBackend}

const Mangai = "Mangai"
