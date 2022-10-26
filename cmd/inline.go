package cmd

import (
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/metafates/mangal/anilist"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/converter"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/inline"
	"github.com/metafates/mangal/provider"
	"github.com/metafates/mangal/util"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
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
	lo.Must0(viper.BindPFlag(constant.MetadataFetchAnilist, inlineCmd.Flags().Lookup("fetch-metadata")))

	inlineCmd.Flags().StringP("output", "o", "", "output file")

	lo.Must0(inlineCmd.MarkFlagRequired("query"))
	inlineCmd.MarkFlagsMutuallyExclusive("download", "json")
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

		if _, err := converter.Get(viper.GetString(constant.FormatsUse)); err != nil {
			handleErr(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		sourceName := viper.GetString(constant.DownloaderDefaultSource)
		if sourceName == "" {
			handleErr(errors.New("source not set"))
		}
		p, ok := provider.Get(sourceName)
		if !ok {
			handleErr(fmt.Errorf("source not found: %s", sourceName))
		}

		src, err := p.CreateSource()
		handleErr(err)

		output := lo.Must(cmd.Flags().GetString("output"))
		var writer io.Writer
		if output != "" {
			writer, err = filesystem.Api().Create(output)
			handleErr(err)
		} else {
			writer = os.Stdout
		}

		mangaFlag := lo.Must(cmd.Flags().GetString("manga"))
		mangaPicker := util.None[inline.MangaPicker]()
		if mangaFlag != "" {
			fn, err := inline.ParseMangaPicker(lo.Must(cmd.Flags().GetString("manga")))
			handleErr(err)
			mangaPicker = util.Some(fn)
		}

		chapterFlag := lo.Must(cmd.Flags().GetString("chapters"))
		chapterFilter := util.None[inline.ChaptersFilter]()
		if chapterFlag != "" {
			fn, err := inline.ParseChaptersFilter(chapterFlag)
			handleErr(err)
			chapterFilter = util.Some(fn)
		}

		options := &inline.Options{
			Source:         src,
			Download:       lo.Must(cmd.Flags().GetBool("download")),
			Json:           lo.Must(cmd.Flags().GetBool("json")),
			Query:          lo.Must(cmd.Flags().GetString("query")),
			PopulatePages:  lo.Must(cmd.Flags().GetBool("populate-pages")),
			MangaPicker:    mangaPicker,
			ChaptersFilter: chapterFilter,
			Out:            writer,
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
		name := lo.Must(cmd.Flags().GetString("name"))
		anilistManga, ok := anilist.GetRelation(name)

		if !ok {
			var err error
			anilistManga, err = anilist.FindClosest(name)
			handleErr(err)
		}

		handleErr(json.NewEncoder(os.Stdout).Encode(anilistManga))
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
