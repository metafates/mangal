package info

import (
	"io"
	"net/url"
	"strings"
	"text/template"

	"github.com/mangalorg/libmangal"
	"github.com/pelletier/go-toml"
	"github.com/samber/lo"
)

//go:generate enumer -type=Type -trimprefix=Type -json -text
type Type uint8

const (
	TypeBundle Type = iota + 1
	TypeLua
)

const Filename = "mangal.toml"

// Info contains libmangal info about provider with mangal specific type field
type Info struct {
	libmangal.ProviderInfo
	Type Type `json:"type"`
}

// New parses info from reader
func New(r io.Reader) (info Info, err error) {
	decoder := toml.NewDecoder(r)
	decoder.Strict(true)
	decoder.SetTagName("json")

	err = decoder.Decode(&info)
	return
}

func (i Info) Markdown() string {
	tmpl := template.Must(template.New("markdown").Funcs(map[string]any{
		"domain": func(URLString string) string {
			URL := lo.Must(url.Parse(URLString))
			return URL.Hostname()
		},
	}).Parse(`
{{ with .Provider }}
# {{ .Name }} v{{ .Version }}

{{ if .Website }}
Mangal provider for the [{{ domain .Website }}]({{ .Website }})
{{ end }}

> {{ .Description }}
{{ end }}
`))

	var sb strings.Builder
	lo.Must0(tmpl.Execute(&sb, i))

	return strings.TrimSpace(sb.String())
}
