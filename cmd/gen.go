package cmd

import (
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/util"
	"github.com/metafates/mangal/where"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/user"
	"path/filepath"
	"strings"
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
			"max": func(nums ...int) int {
				max := nums[0]
				for _, num := range nums {
					if num > max {
						max = num
					}
				}

				return max
			},
		}

		tmpl, err := template.New("source").Funcs(funcMap).Parse(constant.SourceTemplate)
		handleErr(err)

		target := filepath.Join(where.Sources(), util.SanitizeFilename(s.Name)+".lua")
		f, err := filesystem.Get().Create(target)
		handleErr(err)

		util.Ignore(f.Close)

		err = tmpl.Execute(f, s)
		handleErr(err)

		cmd.Println(target)
	},
}
