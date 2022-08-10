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
	"github.com/spf13/viper"
)

var highlightCommand = cobra.Command{
	Use:   "highlight ID",
	Short: "scrape highlight",
	Long:  "scrape highlight",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.Unmarshal(&conf.Configuration)
	},
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("highlight expects a username as the first argument")
		}
		return nil
	},
	RunE: func(_ *cobra.Command, args []string) error {
		highlightID := args[0]
		instagram := shared.NewInstagram(conf.Configuration.FBSR, conf.Configuration.Session, conf.Configuration.User)
		URLs, username, err := instagram.Reels(highlightID, true)
		if err != nil {
			return err
		}

		log.Printf("found %d files\n", len(URLs))
		fileNames := make([]string, 0, len(URLs))

		for _, urlString := range URLs {
			URL, parsingError := url.Parse(urlString)
			if parsingError != nil {
				err = fmt.Errorf("%v\n%v", err, parsingError)
				continue
			}

			fileName := fmt.Sprintf("%s_%s_%s_%s", types.Highlight, username, highlightID, path.Base(URL.Path))
			fileNames = append(fileNames, fileName)
		}

		if errs := conf.SaveBundle(types.Highlight, fileNames, URLs); len(errs) != 0 {
			for _, saveError := range errs {
				err = fmt.Errorf("%v\n%v", err, saveError)
			}
		}

		return err
	},
}

func init() {
	RootCommand.AddCommand(&highlightCommand)
}
