package mustache

import (
	"io/ioutil"
	"os"
	"path"
	"regexp"
)

// PartialProvider comprises the behaviors required of a struct to be able to provide partials to the mustache rendering
// engine.
type PartialProvider interface {
	// Get accepts the name of a partial and returns the parsed partial, if it could be found; a valid but empty
	// template, if it could not be found; or nil and error if an error occurred (other than an inability to find
	// the partial).
	Get(name string) (string, error)
}

// FileProvider implements the PartialProvider interface by providing partials drawn from a filesystem. When a partial
// named `NAME`  is requested, FileProvider searches each listed path for a file named as `NAME` followed by any of the
// listed extensions. The default for `Paths` is to search the current working directory. The default for `Extensions`
// is to examine, in order, no extension; then ".mustache"; then ".stache".
type FileProvider struct {
	Paths      []string
	Extensions []string
}

// Get accepts the name of a partial and returns the parsed partial.
func (fp *FileProvider) Get(name string) (string, error) {
	var filename string

	var paths []string
	if fp.Paths != nil {
		paths = fp.Paths
	} else {
		paths = []string{""}
	}

	var exts []string
	if fp.Extensions != nil {
		exts = fp.Extensions
	} else {
		exts = []string{"", ".mustache", ".stache"}
	}

	for _, p := range paths {
		for _, e := range exts {
			name := path.Join(p, name+e)
			f, err := os.Open(name)
			if err == nil {
				filename = name
				f.Close()
				break
			}
		}
	}

	if filename == "" {
		return "", nil
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

var _ PartialProvider = (*FileProvider)(nil)

// StaticProvider implements the PartialProvider interface by providing partials drawn from a map, which maps partial
// name to template contents.
type StaticProvider struct {
	Partials map[string]string
}

// Get accepts the name of a partial and returns the parsed partial.
func (sp *StaticProvider) Get(name string) (string, error) {
	if sp.Partials != nil {
		if data, ok := sp.Partials[name]; ok {
			return data, nil
		}
	}

	return "", nil
}

var _ PartialProvider = (*StaticProvider)(nil)

func getPartials(partials PartialProvider, name, indent string) (*Template, error) {
	data, err := partials.Get(name)
	if err != nil {
		return nil, err
	}

	// indent non empty lines
	r := regexp.MustCompile(`(?m:^(.+)$)`)
	data = r.ReplaceAllString(data, indent+"$1")

	return ParseStringPartials(data, partials)
}
