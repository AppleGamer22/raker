package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var storyCommand = cobra.Command{
	Use:   "story",
	Short: "scrape story",
	Long:  "scrape story",
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("story expects a post ID as the first argument")
		}
		return nil
	},
	RunE: func(_ *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	RootCommand.AddCommand(&storyCommand)
}
