package cmd

import (
	"encoding/json"
	"strings"
	"text/template"

	"github.com/charmbracelet/lipgloss"
	"github.com/mangalorg/mangal/color"
	"github.com/mangalorg/mangal/nametemplate/util"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

func init() {
	subcommands = append(subcommands, templatesCmd)
}

var templatesCmd = &cobra.Command{
	Use:     "templates",
	Aliases: []string{"t"},
	Short:   "Command related to the name templates",
}

func init() {
	templatesCmd.AddCommand(templatesFuncsCmd)
}

var templatesFuncsCmd = &cobra.Command{
	Use:   "funcs",
	Short: "Show available name template functions",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		keyStyle := lipgloss.NewStyle().Bold(true).Foreground(color.Accent)
		descriptionStyle := lipgloss.NewStyle().Italic(true)

		for k, v := range util.Funcs {
			cmd.Println(keyStyle.Render(k))
			cmd.Println(descriptionStyle.Render(v.Description))
			cmd.Println()
		}
	},
}

var templatesExecArgs = struct {
	Value string
}{}

func init() {
	templatesCmd.AddCommand(templatesExecCmd)

	exampleValue := struct {
		Title  string
		Number float64
	}{
		Title:  "Example Title",
		Number: 32.5,
	}

	marshalled := lo.Must(json.Marshal(exampleValue))

	templatesExecCmd.Flags().StringVarP(&templatesExecArgs.Value, "value", "v", string(marshalled), "JSON object to use as value")
}

var templatesExecCmd = &cobra.Command{
	Use:   "exec template...",
	Short: "Execute template",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tmpl, err := template.
			New("exec").
			Funcs(util.FuncMap).
			Parse(strings.Join(args, " "))

		if err != nil {
			errorf(cmd, err.Error())
		}

		var value map[string]any

		if err := json.Unmarshal([]byte(templatesExecArgs.Value), &value); err != nil {
			errorf(cmd, err.Error())
		}

		if err := tmpl.Execute(cmd.OutOrStdout(), value); err != nil {
			errorf(cmd, err.Error())
		}

		cmd.Println()
	},
}
