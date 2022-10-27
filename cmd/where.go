package cmd

import (
	"github.com/metafates/mangal/color"
	"os"

	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/style"
	"github.com/metafates/mangal/util"
	"github.com/metafates/mangal/where"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var wherePaths = []lo.Tuple2[string, func() string]{
	{"downloads", where.Downloads},
	{"config", where.Config},
	{"sources", where.Sources},
	{"logs", where.Logs},
}

func init() {
	rootCmd.AddCommand(whereCmd)

	for _, n := range wherePaths {
		whereCmd.Flags().BoolP(n.A, string(n.A[0]), false, n.A+" path")
	}

	whereCmd.MarkFlagsMutuallyExclusive(lo.Map(wherePaths, func(t lo.Tuple2[string, func() string], _ int) string {
		return t.A
	})...)

	whereCmd.SetOut(os.Stdout)
}

var whereCmd = &cobra.Command{
	Use:   "where",
	Short: "Show the paths for a files related to the " + constant.Mangal,
	Run: func(cmd *cobra.Command, args []string) {
		headerStyle := style.New().Bold(true).Foreground(color.HiPurple).Render

		for _, n := range wherePaths {
			if lo.Must(cmd.Flags().GetBool(n.A)) {
				cmd.Println(n.B())
				return
			}
		}

		for i, n := range wherePaths {
			cmd.Printf("%s %s\n", headerStyle(util.Capitalize(n.A)+"?"), style.Fg(color.Yellow)("--"+n.A))
			cmd.Println(n.B())

			if i < len(wherePaths)-1 {
				cmd.Println()
			}
		}
	},
}
