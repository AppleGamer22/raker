package commands

import (
	"errors"
	"fmt"
	"net/url"
	"path"

	"github.com/AppleGamer22/raker/cli/conf"
	"github.com/AppleGamer22/raker/shared"
	"github.com/AppleGamer22/raker/shared/types"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var storyCommand = cobra.Command{
	Use:   "story USERNAME",
	Short: "scrape story",
	Long:  "scrape story",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.Unmarshal(&conf.Config)
	},
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("story expects a username as the first argument")
		}
		return nil
	},
	RunE: func(_ *cobra.Command, args []string) error {
		username := args[0]
		instagram := shared.NewInstagram(conf.Config.Instagram.FBSR, conf.Config.Instagram.Session, conf.Config.Instagram.User)
		URLs, username, err := instagram.Reels(username, false)
		if err != nil {
			return err
		}

		log.Debugf("found %d files", len(URLs))
		fileNames := make([]string, 0, len(URLs))

		for _, urlString := range URLs {
			URL, parsingError := url.Parse(urlString)
			if parsingError != nil {
				err = fmt.Errorf("%v\n%v", err, parsingError)
				continue
			}

			fileName := fmt.Sprintf("%s_%s_%s", types.Story, username, path.Base(URL.Path))
			fileNames = append(fileNames, fileName)
		}

		if errs := conf.SaveBundle(types.Story, fileNames, URLs); len(errs) != 0 {
			for _, saveError := range errs {
				err = fmt.Errorf("%v\n%v", err, saveError)
			}
		}

		return err
	},
}

func init() {
	RootCommand.AddCommand(&storyCommand)
}
