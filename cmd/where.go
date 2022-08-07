package cmd

import (
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/style"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(whereCmd)

	whereCmd.Flags().BoolP("config", "c", false, "configuration path")
	whereCmd.Flags().BoolP("sources", "s", false, "sources path")
}

var whereCmd = &cobra.Command{
	Use:   "where",
	Short: "Show the paths for a files related to the mangal",
	Run: func(cmd *cobra.Command, args []string) {
		headerStyle := style.Combined(style.Bold, style.HiBlue)

		whereConfig, _ := cmd.Flags().GetBool("config")
		whereSources, _ := cmd.Flags().GetBool("sources")

		printConfigPath := func() {
			cmd.Println(headerStyle("Configuration path:"))
			for _, path := range lo.Must(config.Paths()) {
				cmd.Println(style.Italic(path))
			}
		}

		printSourcesPath := func() {
			cmd.Println(headerStyle("Sources path:"))
			cmd.Println(style.Italic(viper.GetString(config.SourcesPath)))
		}

		if whereConfig {
			printConfigPath()
		}

		if whereSources {
			printSourcesPath()
		}

		if !whereConfig && !whereSources {
			printConfigPath()

			cmd.Println()

			printSourcesPath()
		}
	},
}
