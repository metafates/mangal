package cmd

import (
	"fmt"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/tui"
	"github.com/metafates/mangal/util"
	"github.com/spf13/viper"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/icon"
	"github.com/metafates/mangal/provider"
	"github.com/metafates/mangal/style"
	"github.com/metafates/mangal/where"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(sourcesCmd)
	sourcesCmd.Flags().BoolP("raw", "r", false, "do not print headers")
	sourcesCmd.SetOut(os.Stdout)
}

var sourcesCmd = &cobra.Command{
	Use:     "sources",
	Short:   "List an available sources",
	Example: "mangal sources",
	Run: func(cmd *cobra.Command, args []string) {
		printHeader := !lo.Must(cmd.Flags().GetBool("raw"))
		headerStyle := style.Combined(style.Bold, style.HiBlue)
		h := func(s string) {
			if printHeader {
				cmd.Println(headerStyle(s))
			}
		}

		defaultProviders := provider.DefaultProviders()
		customProviders := provider.CustomProviders()

		h("Builtin:")
		for _, p := range defaultProviders {
			cmd.Println(p.Name)
		}

		if len(customProviders) == 0 {
			return
		}

		h("")

		h("Custom:")
		for _, p := range provider.CustomProviders() {
			cmd.Println(p.Name)
		}
	},
}

func init() {
	sourcesCmd.AddCommand(sourcesRemoveCmd)
}

var sourcesRemoveCmd = &cobra.Command{
	Use:     "remove",
	Short:   "Remove a custom source",
	Example: "mangal sources remove <name>",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		for _, name := range args {
			path := filepath.Join(where.Sources(), name+provider.CustomProviderExtension)
			handleErr(filesystem.Api().Remove(path))
			fmt.Printf("%s successfully removed %s\n", icon.Get(icon.Success), style.Yellow(name))
		}
	},
}

func init() {
	sourcesCmd.AddCommand(sourcesInstallCmd)
}

var sourcesInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Browse and install custom scrapers",
	Long: `Browse and install custom scrapers from official GitHub repo.
https://github.com/metafates/mangal-scrapers`,
	Run: func(cmd *cobra.Command, args []string) {
		handleErr(tui.Run(&tui.Options{Install: true}))
	},
}

func init() {
	sourcesCmd.AddCommand(sourcesGenCmd)

	sourcesGenCmd.Flags().StringP("name", "n", "", "name of the source")
	sourcesGenCmd.Flags().StringP("url", "u", "", "url of the website")

	lo.Must0(sourcesGenCmd.MarkFlagRequired("name"))
	lo.Must0(sourcesGenCmd.MarkFlagRequired("url"))
}

var sourcesGenCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate a new lua source",
	Long:  `Generate a new lua source.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.SetOut(os.Stdout)

		author := viper.GetString(constant.GenAuthor)
		if author == "" {
			usr, err := user.Current()
			if err == nil {
				author = usr.Username
			} else {
				author = "Anonymous"
			}
		}

		s := struct {
			Name            string
			URL             string
			SearchMangaFn   string
			MangaChaptersFn string
			ChapterPagesFn  string
			Author          string
		}{
			Name:            lo.Must(cmd.Flags().GetString("name")),
			URL:             lo.Must(cmd.Flags().GetString("url")),
			SearchMangaFn:   constant.SearchMangaFn,
			MangaChaptersFn: constant.MangaChaptersFn,
			ChapterPagesFn:  constant.ChapterPagesFn,
			Author:          author,
		}

		funcMap := template.FuncMap{
			"repeat": strings.Repeat,
			"plus":   func(a, b int) int { return a + b },
			"max":    util.Max[int],
		}

		tmpl, err := template.New("source").Funcs(funcMap).Parse(constant.SourceTemplate)
		handleErr(err)

		target := filepath.Join(where.Sources(), util.SanitizeFilename(s.Name)+".lua")
		f, err := filesystem.Api().Create(target)
		handleErr(err)

		defer util.Ignore(f.Close)

		err = tmpl.Execute(f, s)
		handleErr(err)

		cmd.Println(target)
	},
}
