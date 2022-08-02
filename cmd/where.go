package cmd

import (
	"github.com/metafates/mangal/config"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var whereCmd = &cobra.Command{
	Use:   "where",
	Short: "Show the paths to the configuration",
	Run: func(cmd *cobra.Command, args []string) {
		configOnly, _ := cmd.Flags().GetBool("config")
		sourcesOnly, _ := cmd.Flags().GetBool("sources")

		printConfigPath := func() {
			for _, path := range lo.Must(config.Paths()) {
				cmd.Println(path)
			}
		}

		printSourcesPath := func() {
			cmd.Println(viper.GetString(config.SourcesPath))
		}

		if configOnly {
			printConfigPath()
		} else if sourcesOnly {
			printSourcesPath()
		} else {
			cmd.Println("Configuration path:")
			printConfigPath()

			cmd.Println()

			cmd.Println("Sources path:")
			printSourcesPath()
		}
	},
}

func init() {
	rootCmd.AddCommand(whereCmd)

	whereCmd.Flags().BoolP("config", "c", false, "configuration path")
	whereCmd.Flags().BoolP("sources", "s", false, "sources path")

	whereCmd.MarkFlagsMutuallyExclusive("config", "sources")
}
