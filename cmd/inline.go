package cmd

import (
	"fmt"
	"github.com/metafates/mangal/inline"
	"github.com/metafates/mangal/provider"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(inlineCmd)

	inlineCmd.Flags().String("source", "", "source to use. see `mangal sources` for available sources")
	inlineCmd.Flags().String("query", "", "query to search for")
	inlineCmd.Flags().String("manga", "", "manga selector. first|last|[number]|@[substring]@")
	inlineCmd.Flags().String("chapters", "", "chapter selector. first|last|all|[number]|[number]-[number]")
	inlineCmd.Flags().BoolP("download", "d", false, "download chapters")

	lo.Must0(inlineCmd.MarkFlagRequired("source"))
	lo.Must0(inlineCmd.MarkFlagRequired("query"))
	lo.Must0(inlineCmd.MarkFlagRequired("manga"))
	lo.Must0(inlineCmd.MarkFlagRequired("chapters"))
}

var inlineCmd = &cobra.Command{
	Use:     "inline",
	Short:   "Inline mode for scripting",
	Long:    "Inline mode for scripting",
	Example: "mangal inline --source Manganelo --query \"death note\" --manga first --chapters \"@Vol.1 @\" -d",
	RunE: func(cmd *cobra.Command, args []string) error {
		sourceName := lo.Must(cmd.Flags().GetString("source"))
		p, ok := provider.Get(sourceName)
		if !ok {
			return fmt.Errorf("source not found: %s", sourceName)
		}

		src, err := p.CreateSource()
		if err != nil {
			return err
		}

		mangaPicker, err := inline.ParseMangaPicker(lo.Must(cmd.Flags().GetString("manga")))
		if err != nil {
			return err
		}

		chapterFilter, err := inline.ParseChaptersFilter(lo.Must(cmd.Flags().GetString("chapters")))
		if err != nil {
			return err
		}

		options := &inline.Options{
			Source:        src,
			Download:      lo.Must(cmd.Flags().GetBool("download")),
			Query:         lo.Must(cmd.Flags().GetString("query")),
			MangaPicker:   mangaPicker,
			ChapterFilter: chapterFilter,
		}

		return inline.Run(options)
	},
}
