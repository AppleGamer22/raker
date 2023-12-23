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

var tiktokCommand = cobra.Command{
	Use:     "tiktok USERNAME POST",
	Short:   "scrape tiktok",
	Long:    "scrape tiktok",
	Aliases: []string{"tt"},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !incognito {
			return viper.Unmarshal(&conf.Config)
		}
		return nil
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
		tiktok := shared.NewTikTok(conf.Config.TikTok.Session, conf.Config.TikTok.Guard, conf.Config.TikTok.Chain)
		URLs, username, _, err := tiktok.Post(username, post, incognito)
		if err != nil {
			return err
		}
		log.Debug("found %d files", len(URLs))
		fileNames := make([]string, 0, len(URLs))

		for _, urlString := range URLs {
			URL, parsingError := url.Parse(urlString)
			if parsingError != nil {
				err = fmt.Errorf("%v\n%v", err, parsingError)
				continue
			}

			if URL.Query().Get("mime_type") == "video_mp4" {
				fileNames = append(fileNames, fmt.Sprintf("%s.mp4", post))
				break
			}

			fileName := fmt.Sprintf("%s_%s_%s_%s", types.TikTok, username, post, path.Base(URL.Path))
			fileNames = append(fileNames, fileName)
		}

		if errs := conf.SaveBundle(types.TikTok, fileNames, URLs); len(errs) != 0 {
			for _, saveError := range errs {
				err = fmt.Errorf("%v\n%v", err, saveError)
			}
		}
		return err
	},
}

func init() {
	tiktokCommand.Flags().BoolVarP(&incognito, "incognito", "i", false, "without authentication")
	RootCommand.AddCommand(&tiktokCommand)
}
