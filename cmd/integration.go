package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/integration/anilist"
	"github.com/metafates/mangal/log"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(integrationCmd)
	integrationCmd.AddCommand(integrationAnilistCmd)
	integrationAnilistCmd.Flags().BoolP("disable", "d", false, "Disable Anilist integration")
}

var integrationCmd = &cobra.Command{
	Use:   "integration",
	Short: "Integration with other services",
	Long:  `Integration with other services`,
}

var integrationAnilistCmd = &cobra.Command{
	Use:   "anilist",
	Short: "Integration with Anilist",
	Long: `Integration with Anilist.
See https://github.com/metafates/mangal/wiki/Anilist-Integration for more information`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if lo.Must(cmd.Flags().GetBool("disable")) {
			viper.Set(constant.AnilistEnable, false)
			viper.Set(constant.AnilistCode, "")
			viper.Set(constant.AnilistSecret, "")
			viper.Set(constant.AnilistID, "")
			log.Info("Anilist integration disabled")
			return viper.WriteConfig()
		}

		if !viper.GetBool(constant.AnilistEnable) {
			confirm := survey.Confirm{
				Message: "Anilist is disabled. Enable?",
				Default: false,
			}
			var response bool
			err := survey.AskOne(&confirm, &response)
			if err != nil {
				log.Error(err)
				return err
			}

			if !response {
				return nil
			}

			viper.Set(constant.AnilistEnable, response)
			err = viper.WriteConfig()
			if err != nil {
				switch err.(type) {
				case viper.ConfigFileNotFoundError:
					err = viper.SafeWriteConfig()
					if err != nil {
						return err
					}
				default:
					log.Error(err)
					return err
				}
			}
		}

		if viper.GetString(constant.AnilistID) == "" {
			input := survey.Input{
				Message: "Anilsit client ID is not set. Please enter it:",
				Help:    "",
			}
			var response string
			err := survey.AskOne(&input, &response)
			if err != nil {
				return err
			}

			if response == "" {
				return nil
			}

			viper.Set(constant.AnilistID, response)
			err = viper.WriteConfig()
			if err != nil {
				log.Error(err)
				return err
			}
		}

		if viper.GetString(constant.AnilistSecret) == "" {
			input := survey.Input{
				Message: "Anilsit client secret is not set. Please enter it:",
				Help:    "",
			}
			var response string
			err := survey.AskOne(&input, &response)
			if err != nil {
				return err
			}

			if response == "" {
				return nil
			}

			viper.Set(constant.AnilistSecret, response)
			err = viper.WriteConfig()
			if err != nil {
				log.Error(err)
				return err
			}
		}

		if viper.GetString(constant.AnilistCode) == "" {
			fmt.Println(anilist.New().AuthURL())
			input := survey.Input{
				Message: "Anilsit code is not set. Please copy it from the link above and paste in here:",
				Help:    "",
			}
			var response string
			err := survey.AskOne(&input, &response)
			if err != nil {
				return err
			}

			if response == "" {
				return nil
			}

			viper.Set(constant.AnilistCode, response)
			err = viper.WriteConfig()
			if err != nil {
				log.Error(err)
				return err
			}
		}

		return nil
	},
}
