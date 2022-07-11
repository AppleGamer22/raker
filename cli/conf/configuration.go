package conf

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type configuration struct {
	Session string
	FBSR    string
	User    string
	TikTok  string
}

var Configuration configuration

func init() {
	viper.AutomaticEnv()
	viper.BindEnv("FBSR")
	viper.BindEnv("SESSION")
	viper.BindEnv("USER")
	viper.BindEnv("TIKTOK")
	viper.SetConfigName(".rake")
	viper.SetConfigType("yaml")

	directory, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	viper.AddConfigPath(directory)

	viper.ReadInConfig()
}
