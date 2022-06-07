package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AppleGamer22/rake/shared"
	"github.com/spf13/cobra"
)

var instagramSignInCommand = cobra.Command{
	Use:   "in",
	Short: "instagram sign-in",
	Long:  "instagram sign-in",
	RunE: func(_ *cobra.Command, args []string) error {
		var username string
		fmt.Print("username: ")
		if _, err := fmt.Scan(&username); err != nil {
			return err
		}

		var password string
		fmt.Print("password: ")
		if err := readPassword(&password); err != nil {
			return err
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		userDataDir := filepath.Join(homeDir, ".rake")
		raker, err := shared.NewRaker(userDataDir, debug, incognito)
		if err != nil {
			return err
		}
		return raker.InstagramSignIn(username, password)
	},
}

var instagramSignOutCommand = cobra.Command{
	Use:   "out",
	Short: "instagram sign-out",
	Long:  "instagram sign-out",
	RunE: func(_ *cobra.Command, args []string) error {
		var username string
		fmt.Print("username: ")
		if _, err := fmt.Scan(&username); err != nil {
			return err
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		userDataDir := filepath.Join(homeDir, ".rake")
		raker, err := shared.NewRaker(userDataDir, debug, incognito)
		if err != nil {
			return err
		}

		return raker.InstagramSignOut(username)
	},
}

func init() {
	instagramSignInCommand.Flags().BoolVarP(&debug, "debug", "d", false, "visible browser")
	instagramSignInCommand.Flags().BoolVarP(&incognito, "incognito", "i", false, "incognito browser")
	instagramSignOutCommand.Flags().BoolVarP(&debug, "debug", "d", false, "visible browser")
	instagramSignOutCommand.Flags().BoolVarP(&incognito, "incognito", "i", false, "incognito browser")
	instagramCommand.AddCommand(&instagramSignInCommand)
	instagramCommand.AddCommand(&instagramSignOutCommand)
}
