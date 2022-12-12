package shared_test

import (
	"github.com/AppleGamer22/raker/cli/conf"
	"github.com/spf13/viper"
)

var configuration conf.Configuration

func init() {
	// https://stackoverflow.com/a/65747120/7148921
	viper.AutomaticEnv()
	viper.BindEnv("instagram.session", "SESSION_IG")
	viper.BindEnv("instagram.user", "USER")
	viper.BindEnv("instagram.fbsr", "FBSR")
	viper.BindEnv("tiktok.session", "SESSION_TT")
	viper.BindEnv("tiktok.chain", "TIKTOK_CT")
	viper.SetConfigName(".raker")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("..")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&configuration); err != nil {
		panic(err)
	}
}
