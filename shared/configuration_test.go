package shared_test

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	Instagram struct {
		Session string
		FBSR    string
		App     string
	}
	TikTok string
}

var configuration Configuration

func init() {
	viper.SetConfigName(".rake")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("..")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&configuration); err != nil {
		panic(err)
	}
}
