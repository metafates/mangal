package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/invopop/jsonschema"
	"github.com/metafates/mangal/anilist"
	"github.com/metafates/mangal/converter"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/inline"
	"github.com/metafates/mangal/key"
	"github.com/metafates/mangal/provider"
	"github.com/metafates/mangal/query"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/update"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func init() {
	rootCmd.AddCommand(inlineCmd)

	inlineCmd.Flags().StringP("query", "q", "", "query to search for")
	inlineCmd.Flags().StringP("manga", "m", "", "manga selector")
	inlineCmd.Flags().StringP("chapters", "c", "", "chapter selector")
	inlineCmd.Flags().BoolP("download", "d", false, "download chapters")
	inlineCmd.Flags().BoolP("json", "j", false, "JSON output")
	inlineCmd.Flags().BoolP("populate-pages", "p", false, "Populate chapters pages")
	inlineCmd.Flags().BoolP("fetch-metadata", "f", false, "Populate manga metadata")
	inlineCmd.Flags().BoolP("include-anilist-manga", "a", false, "Include anilist manga in the output")
	lo.Must0(viper.BindPFlag(key.MetadataFetchAnilist, inlineCmd.Flags().Lookup("fetch-metadata")))

	inlineCmd.Flags().StringP("output", "o", "", "output file")

	lo.Must0(inlineCmd.MarkFlagRequired("query"))
	inlineCmd.MarkFlagsMutuallyExclusive("download", "json")
	inlineCmd.MarkFlagsMutuallyExclusive("include-anilist-manga", "download")

	inlineCmd.RegisterFlagCompletionFunc("query", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return query.SuggestMany(toComplete), cobra.ShellCompDirectiveNoFileComp
	})
}

var inlineCmd = &cobra.Command{
	Use:   "inline",
	Short: "Launch in the inline mode",
	Long: `Launch in the inline mode for scripting

Manga selectors:
  first - first manga in the list
  last - last manga in the list
  [number] - select manga by index (starting from 0)

Chapter selectors:
  first - first chapter in the list
  last - last chapter in the list
  all - all chapters in the list
  [number] - select chapter by index (starting from 0)
  [from]-[to] - select chapters by range
  @[substring]@ - select chapters by name substring

When using the json flag manga selector could be omitted. That way, it will select all mangas`,

	Example: "https://github.com/metafates/mangal/wiki/Inline-mode",
	PreRun: func(cmd *cobra.Command, args []string) {
		json, _ := cmd.Flags().GetBool("json")

		if !json {
			lo.Must0(cmd.MarkFlagRequired("manga"))
		}

		if lo.Must(cmd.Flags().GetBool("populate-pages")) {
			lo.Must0(cmd.MarkFlagRequired("json"))
		}

		if _, err := converter.Get(viper.GetString(key.FormatsUse)); err != nil {
			handleErr(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var (
			sources []source.Source
			err     error
		)

		for _, name := range viper.GetStringSlice(key.DownloaderDefaultSources) {
			if name == "" {
				handleErr(errors.New("source not set"))
			}

			p, ok := provider.Get(name)
			if !ok {
				handleErr(fmt.Errorf("source not found: %s", name))
			}

			src, err := p.CreateSource()
			handleErr(err)

			sources = append(sources, src)
		}

		query := lo.Must(cmd.Flags().GetString("query"))

		output := lo.Must(cmd.Flags().GetString("output"))
		var writer io.Writer
		if output != "" {
			writer, err = filesystem.Api().Create(output)
			handleErr(err)
		} else {
			writer = os.Stdout
		}

		mangaFlag := lo.Must(cmd.Flags().GetString("manga"))
		mangaPicker := mo.None[inline.MangaPicker]()
		if mangaFlag != "" {
			fn, err := inline.ParseMangaPicker(query, lo.Must(cmd.Flags().GetString("manga")))
			handleErr(err)
			mangaPicker = mo.Some(fn)
		}

		chapterFlag := lo.Must(cmd.Flags().GetString("chapters"))
		chapterFilter := mo.None[inline.ChaptersFilter]()
		if chapterFlag != "" {
			fn, err := inline.ParseChaptersFilter(chapterFlag)
			handleErr(err)
			chapterFilter = mo.Some(fn)
		}

		options := &inline.Options{
			Sources:             sources,
			Download:            lo.Must(cmd.Flags().GetBool("download")),
			Json:                lo.Must(cmd.Flags().GetBool("json")),
			Query:               query,
			PopulatePages:       lo.Must(cmd.Flags().GetBool("populate-pages")),
			IncludeAnilistManga: lo.Must(cmd.Flags().GetBool("include-anilist-manga")),
			MangaPicker:         mangaPicker,
			ChaptersFilter:      chapterFilter,
			Out:                 writer,
		}

		handleErr(inline.Run(options))
	},
}

func init() {
	inlineCmd.AddCommand(inlineAnilistCmd)
}

var inlineAnilistCmd = &cobra.Command{
	Use:   "anilist",
	Short: "Anilist related commands",
}

func init() {
	inlineAnilistCmd.AddCommand(inlineAnilistSearchCmd)

	inlineAnilistSearchCmd.Flags().StringP("name", "n", "", "manga name to search")
	inlineAnilistSearchCmd.Flags().IntP("id", "i", 0, "anilist manga id")

	inlineAnilistSearchCmd.MarkFlagsMutuallyExclusive("name", "id")
}

var inlineAnilistSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search anilist manga by name",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("name") && !cmd.Flags().Changed("id") {
			handleErr(errors.New("name or id flag is required"))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		mangaName := lo.Must(cmd.Flags().GetString("name"))
		mangaId := lo.Must(cmd.Flags().GetInt("id"))

		var toEncode any

		if mangaName != "" {
			mangas, err := anilist.SearchByName(mangaName)
			handleErr(err)
			toEncode = mangas
		} else {
			manga, err := anilist.GetByID(mangaId)
			handleErr(err)
			toEncode = manga
		}

		handleErr(json.NewEncoder(os.Stdout).Encode(toEncode))
	},
}

func init() {
	inlineAnilistCmd.AddCommand(inlineAnilistGetCmd)

	inlineAnilistGetCmd.Flags().StringP("name", "n", "", "manga name to get the bind for")
	lo.Must0(inlineAnilistGetCmd.MarkFlagRequired("name"))
}

var inlineAnilistGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get anilist manga that is bind to manga name",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			m   *anilist.Manga
			err error
		)

		name := lo.Must(cmd.Flags().GetString("name"))
		m, err = anilist.FindClosest(name)

		if err != nil {
			m, err = anilist.FindClosest(name)
			handleErr(err)
		}

		handleErr(json.NewEncoder(os.Stdout).Encode(m))
	},
}

func init() {
	inlineAnilistCmd.AddCommand(inlineAnilistBindCmd)

	inlineAnilistBindCmd.Flags().StringP("name", "n", "", "manga name")
	inlineAnilistBindCmd.Flags().IntP("id", "i", 0, "anilist manga id")

	lo.Must0(inlineAnilistBindCmd.MarkFlagRequired("name"))
	lo.Must0(inlineAnilistBindCmd.MarkFlagRequired("id"))

	inlineAnilistBindCmd.MarkFlagsRequiredTogether("name", "id")
}

var inlineAnilistBindCmd = &cobra.Command{
	Use:   "set",
	Short: "Bind manga name to the anilist manga by id",
	Run: func(cmd *cobra.Command, args []string) {
		anilistManga, err := anilist.GetByID(lo.Must(cmd.Flags().GetInt("id")))
		handleErr(err)

		mangaName := lo.Must(cmd.Flags().GetString("name"))

		handleErr(anilist.SetRelation(mangaName, anilistManga))
	},
}

func init() {
	inlineAnilistCmd.AddCommand(inlineAnilistUpdateCmd)

	inlineAnilistUpdateCmd.Flags().StringP("path", "p", "", "path to the manga")
	lo.Must0(inlineAnilistUpdateCmd.MarkFlagRequired("path"))
}

var inlineAnilistUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update old manga metadata according to the current anilist bind",
	Run: func(cmd *cobra.Command, args []string) {
		path := lo.Must(cmd.Flags().GetString("path"))
		handleErr(update.Metadata(path))
	},
}

func init() {
	inlineCmd.AddCommand(inlineSchemaCmd)

	inlineSchemaCmd.Flags().BoolP("anilist", "a", false, "generate anilist search output schema")
}

var inlineSchemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Schemas for the inline json outputs",
	Run: func(cmd *cobra.Command, args []string) {
		reflector := new(jsonschema.Reflector)
		reflector.Anonymous = true
		reflector.Namer = func(t reflect.Type) string {
			name := t.Name()
			switch strings.ToLower(name) {
			case "manga", "chapter", "page", "date", "output":
				return filepath.Base(t.PkgPath()) + "." + name
			}

			return name
		}

		var schema *jsonschema.Schema

		switch {
		case lo.Must(cmd.Flags().GetBool("anilist")):
			schema = reflector.Reflect([]*anilist.Manga{})
		default:
			schema = reflector.Reflect(&inline.Output{})
		}

		handleErr(json.NewEncoder(os.Stdout).Encode(schema))
	},
}
