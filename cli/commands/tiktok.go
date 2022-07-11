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

var tiktokCommand = cobra.Command{
	Use:     "tiktok USERNAME POST",
	Short:   "scrape tiktok",
	Long:    "scrape tiktok",
	Aliases: []string{"tt"},
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.Unmarshal(&conf.Configuration)
	},
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("tiktok expects a username & post ID as the first argument")
		}
		return nil
	},
	RunE: func(_ *cobra.Command, args []string) error {
		username := args[0]
		post := args[1]
		tiktok := shared.NewTikTok(conf.Configuration.TikTok)
		urlString, username, err := tiktok.Post(username, post)
		if err != nil {
			return err
		}
		log.Println("found 1 file")
		URL, err := url.Parse(urlString)
		if err != nil {
			return err
		}
		fileName := fmt.Sprintf("%s_%s_%s_%s", types.TikTok, username, post, path.Base(URL.Path))
		if err := conf.Save(types.TikTok, fileName, urlString); err != nil {
			return err
		}
		log.Printf("saved %s to file %s at the current directory", urlString, fileName)
		return nil
	},
}

func init() {
	RootCommand.AddCommand(&tiktokCommand)
}
