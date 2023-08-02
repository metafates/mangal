package cmd

import (
	"fmt"
	"github.com/mangalorg/mangal/config"
	"github.com/spf13/cobra"
	"strings"
	"text/template"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
}

func init() {
	configCmd.AddCommand(configInfoCmd)
}

var configInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show configuration information",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
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
	},
}

func init() {
	configCmd.AddCommand(configWriteCmd)
}

var configWriteCmd = &cobra.Command{
	Use:   "write",
	Short: "Write configuration to disk",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return config.Write()
	},
}
