package cmd

import (
	"github.com/AppleGamer22/rake/shared"
	"github.com/spf13/cobra"
)

var RootCommand = cobra.Command{
	Use:     "rake",
	Short:   "scrape common social media networks",
	Long:    "scrape common social media networks",
	Version: shared.Version,
}
