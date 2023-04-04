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

var instagramCommand = cobra.Command{
	Use:     "instagram POST",
	Short:   "scrape instagram",
	Long:    "scrape instagram",
	Aliases: []string{"ig"},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !incognito {
			return viper.Unmarshal(&conf.Config)
		}
		return nil
	},
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("instagram expects a post ID as the first argument")
		}
		return nil
	},
	RunE: func(_ *cobra.Command, args []string) error {
		post := args[0]
		instagram := shared.NewInstagram(conf.Config.Instagram.FBSR, conf.Config.Instagram.Session, conf.Config.Instagram.User)
		URLs, username, err := instagram.Post(post, incognito)
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

			fileName := fmt.Sprintf("%s_%s_%s_%s", types.Instagram, username, post, path.Base(URL.Path))
			fileNames = append(fileNames, fileName)
		}

		if errs := conf.SaveBundle(types.Instagram, fileNames, URLs); len(errs) != 0 {
			for _, saveError := range errs {
				err = fmt.Errorf("%v\n%v", err, saveError)
			}
		}
		return err
	},
}

func init() {
	instagramCommand.Flags().BoolVarP(&incognito, "incognito", "i", false, "without authentication")
	RootCommand.AddCommand(&instagramCommand)
}
