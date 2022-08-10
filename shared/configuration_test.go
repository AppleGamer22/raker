package shared_test

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	Instagram struct {
		Session string
		User    string
		FBSR    string
	}
	TikTok string
}

var configuration Configuration

func init() {
	// https://stackoverflow.com/a/65747120/7148921
	viper.AutomaticEnv()
	viper.BindEnv("instagram.session", "SESSION")
	viper.BindEnv("instagram.user", "USER")
	viper.BindEnv("instagram.fbsr", "FBSR")
	viper.BindEnv("TIKTOK")
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
