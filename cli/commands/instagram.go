package commands

import (
	"errors"

	"github.com/spf13/cobra"
)

var instagramCommand = cobra.Command{
	Use:     "instagram",
	Short:   "scrape instagram",
	Long:    "scrape instagram",
	Aliases: []string{"ig"},
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("instagram expects a post ID as the first argument")
		}
		return nil
	},
	RunE: func(_ *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	RootCommand.AddCommand(&instagramCommand)
}
