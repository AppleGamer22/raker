package commands

import (
	"github.com/AppleGamer22/raker/shared"
	"github.com/spf13/cobra"
)

var (
	incognito   bool
	RootCommand = cobra.Command{
		Use:     "raker",
		Short:   "scrape common social media networks",
		Long:    "scrape common social media networks",
		Version: shared.Version,
	}
)

func init() {
	RootCommand.SetVersionTemplate("{{.Version}}\n")
}
