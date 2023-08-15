package cmd

import (
	"bufio"
	"context"
	"fmt"
	"strings"

	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/anilist"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

func init() {
	subcommands = append(subcommands, anilistCmd)
}

var anilistCmd = &cobra.Command{
	Use:     "anilist",
	Aliases: []string{"al"},
	Short:   "Anilist related commands",
}

func init() {
	anilistCmd.AddCommand(anilistAuthCmd)
}

var anilistAuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authorize with anilist to enable sync",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: replace it with something better
		reader := bufio.NewReader(cmd.InOrStdin())

		cmd.Print("Enter your ID\n> ")
		id, err := reader.ReadString('\n')
		if err != nil {
			errorf(cmd, err.Error())
		}
		id = strings.TrimSpace(id)

		cmd.Print("Enter your secret\n> ")
		secret, err := reader.ReadString('\n')
		if err != nil {
			errorf(cmd, err.Error())
		}
		secret = strings.TrimSpace(secret)

		authURL := fmt.Sprint("https://anilist.co/api/v2/oauth/authorize?client_id=", id, "&response_type=code&redirect_uri=https://anilist.co/api/v2/oauth/pin")
		if err := open.Start(authURL); err != nil {
			errorf(cmd, err.Error())
		}

		cmd.Print("Enter your code\n> ")
		code, err := reader.ReadString('\n')
		if err != nil {
			errorf(cmd, err.Error())
		}
		code = strings.TrimSpace(code)

		err = anilist.Client.Authorize(context.Background(), libmangal.AnilistLoginCredentials{
			ID:     id,
			Secret: secret,
			Code:   code,
		})

		if err != nil {
			errorf(cmd, err.Error())
		}

		successf(cmd, "Authorized with the Anilist")
	},
}

func init() {
	anilistCmd.AddCommand(anilistLogoutCmd)
}

var anilistLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Exit from anilist",
	Run: func(cmd *cobra.Command, args []string) {
		if err := anilist.Client.Logout(); err != nil {
			errorf(cmd, err.Error())
		}

		successf(cmd, "Logged out from Anilist")
	},
}
