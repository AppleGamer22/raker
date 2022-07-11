package commands

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"path"

	"github.com/AppleGamer22/rake/cli/conf"
	"github.com/AppleGamer22/rake/shared"
	"github.com/AppleGamer22/rake/shared/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var storyCommand = cobra.Command{
	Use:   "story",
	Short: "scrape story",
	Long:  "scrape story",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.Unmarshal(&conf.Configuration)
	},
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("story expects a username as the first argument")
		}
		return nil
	},
	RunE: func(_ *cobra.Command, args []string) error {
		errs := []error{}
		username := args[0]
		instagram := shared.NewInstagram(conf.Configuration.FBSR, conf.Configuration.Session, conf.Configuration.User)
		URLs, username, err := instagram.Reels(username, false)
		if err != nil {
			return err
		}
		log.Printf("found %d files\n", len(URLs))
		for _, urlString := range URLs {
			URL, err := url.Parse(urlString)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			fileName := fmt.Sprintf("%s_%s_%s", types.Story, username, path.Base(URL.Path))
			if err = conf.Save(types.Story, fileName, urlString); err != nil {
				errs = append(errs, err)
				continue
			}
			log.Printf("saved %s to file %s at the current directory\n", urlString, fileName)
		}
		for _, err2 := range errs {
			err = fmt.Errorf("%v\n%v", err, err2)
		}
		return err
	},
}

func init() {
	RootCommand.AddCommand(&storyCommand)
}
