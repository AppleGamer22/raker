package commands

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"path"

	"github.com/AppleGamer22/raker/cli/conf"
	"github.com/AppleGamer22/raker/shared"
	"github.com/AppleGamer22/raker/shared/types"
	"github.com/spf13/cobra"
)

var vscoCommand = cobra.Command{
	Use:   "vsco USERNAME POST",
	Short: "scrape vsco",
	Long:  "scrape vsco",
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("vsco expects a username & post ID as the first argument")
		}
		return nil
	},
	RunE: func(_ *cobra.Command, args []string) error {
		username := args[0]
		post := args[1]
		urlString, username, err := shared.VSCO(username, post)
		if err != nil {
			return err
		}
		log.Println("found 1 file")
		URL, err := url.Parse(urlString)
		if err != nil {
			return err
		}
		fileName := fmt.Sprintf("%s_%s_%s_%s", types.VSCO, username, post, path.Base(URL.Path))
		if err := conf.Save(types.VSCO, fileName, urlString); err != nil {
			return err
		}
		log.Printf("saved %s to file %s at the current directory", urlString, fileName)
		return nil
	},
}

func init() {
	RootCommand.AddCommand(&vscoCommand)
}
