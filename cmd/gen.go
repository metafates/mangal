package cmd

import (
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/util"
	"github.com/metafates/mangal/where"
	"github.com/samber/lo"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"text/template"
)

func init() {
	rootCmd.AddCommand(genCmd)

	genCmd.Flags().StringP("name", "n", "", "name of the source")
	genCmd.Flags().StringP("url", "u", "", "url of the website")

	lo.Must0(genCmd.MarkFlagRequired("name"))
	lo.Must0(genCmd.MarkFlagRequired("url"))
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate a new lua source",
	Long:  `Generate a new lua source.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SetOut(os.Stdout)

		s := struct {
			Name            string
			URL             string
			SearchMangaFn   string
			MangaChaptersFn string
			ChapterPagesFn  string
		}{
			Name:            lo.Must(cmd.Flags().GetString("name")),
			URL:             lo.Must(cmd.Flags().GetString("url")),
			SearchMangaFn:   constant.SearchMangaFn,
			MangaChaptersFn: constant.MangaChaptersFn,
			ChapterPagesFn:  constant.ChapterPagesFn,
		}

		tmpl, err := template.New("source").Parse(constant.SourceTemplate)

		if err != nil {
			return err
		}

		target := filepath.Join(where.Sources(), util.SanitizeFilename(s.Name)+".lua")
		f, err := filesystem.Get().Create(target)
		if err != nil {
			return err
		}

		defer func(f afero.File) {
			_ = f.Close()
		}(f)

		err = tmpl.Execute(f, s)
		if err != nil {
			return err
		}

		cmd.Println(target)
		return nil
	},
}
