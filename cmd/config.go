package cmd

import (
	"fmt"
	"github.com/mangalorg/mangal/config"
	"strings"
	"text/template"
)

type configCmd struct {
	Info configInfoCmd `cmd:""`
}

type configInfoCmd struct{}

func (c *configInfoCmd) Run() error {
	fieldTemplate := template.Must(template.New("field").Parse(`
{{.Description}}

Key: {{.Key}}
Value: {{.Value}}
Default: {{.Default}}
`))

	var sb strings.Builder
	for _, field := range config.Fields {
		if err := fieldTemplate.Execute(&sb, field); err != nil {
			return err
		}

		fmt.Print(sb.String())
		sb.Reset()
	}

	return nil
}
