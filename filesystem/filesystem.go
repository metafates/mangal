package filesystem

import "github.com/spf13/afero"

var appFs = afero.NewOsFs()

// Get returns afero filesystem abstraction layer
func Get() afero.Fs {
	return appFs
}

// Set sets afero filesystem abstraction layer
func Set(fs afero.Fs) {
	appFs = fs
}
