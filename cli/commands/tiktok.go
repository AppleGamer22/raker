package commands

import (
	"errors"

	"github.com/spf13/cobra"
)

var tiktokCommand = cobra.Command{
	Use:   "tiktok",
	Short: "scrape tiktok",
	Long:  "scrape tiktok",
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("tiktok expects a post ID as the first argument")
		}
		return nil
	},
	RunE: func(_ *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	RootCommand.AddCommand(&tiktokCommand)
}
