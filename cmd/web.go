package cmd

import (
	"time"

	"github.com/mangalorg/mangal/web"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(webCmd)

	webCmd.Flags().BoolVarP(&webArgs.Open, "open", "o", false, "Open served page in the default browser")
	webCmd.Flags().StringVarP(&webArgs.Port, "port", "p", "6969", "Port to use")
}

var webArgs = struct {
	Open bool
	Port string
}{}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Run web UI",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if webArgs.Open {
			time.AfterFunc(time.Second, func() {
				open.Run("http://localhost:" + webArgs.Port)
			})
		}
		if err := web.Run(); err != nil {
			errorf(cmd, err.Error())
		}
	},
}
