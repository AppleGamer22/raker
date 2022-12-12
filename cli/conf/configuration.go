package conf

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Configuration struct {
	Instagram struct {
		Session string
		User    string
		FBSR    string
	}
	TikTok struct {
		Session string
		Guard   string
		Chain   string
	}
}

var Config Configuration

func init() {
	viper.AutomaticEnv()
	viper.BindEnv("instagram.session", "SESSION_IG")
	viper.BindEnv("instagram.user", "USER")
	viper.BindEnv("instagram.fbsr", "FBSR")
	viper.BindEnv("tiktok.session", "SESSION_TT")
	viper.BindEnv("tiktok.chain", "TIKTOK_CT")
	viper.BindEnv("tiktok.guard", "GUARD")
	viper.SetConfigName(".raker")
	viper.SetConfigType("yaml")

	directory, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	viper.AddConfigPath(directory)

	viper.ReadInConfig()
}
