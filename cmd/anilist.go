package cmd

import (
	"bufio"
	"context"
	"fmt"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/anilist"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func init() {
	rootCmd.AddCommand(anilistCmd)
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
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: replace it with something better
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter your ID\n> ")
		id, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		id = strings.TrimSpace(id)

		fmt.Print("Enter your secret\n> ")
		secret, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		secret = strings.TrimSpace(secret)

		authURL := fmt.Sprint("https://anilist.co/api/v2/oauth/authorize?client_id=", id, "&response_type=code&redirect_uri=https://anilist.co/api/v2/oauth/pin")
		if err := open.Start(authURL); err != nil {
			return err
		}

		fmt.Print("Enter your code\n> ")
		code, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		code = strings.TrimSpace(code)

		return anilist.Client.Authorize(context.Background(), libmangal.AnilistLoginCredentials{
			ID:     id,
			Secret: secret,
			Code:   code,
		})
	},
}
