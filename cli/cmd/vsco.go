package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var vscoCommand = cobra.Command{
	Use:   "vsco",
	Short: "scrape vsco",
	Long:  "scrape vsco",
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("vsco expects a post ID as the first argument")
		}
		return nil
	},
	RunE: func(_ *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	RootCommand.AddCommand(&vscoCommand)
}
