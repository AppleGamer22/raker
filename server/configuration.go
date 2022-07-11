package main

import "github.com/spf13/viper"

type Configuration struct {
	Secret      string
	URI         string
	Database    string
	Storage     string
	Directories bool
	Port        uint
}

var configuration = Configuration{
	URI:         "mongodb://localhost:27017",
	Database:    "rake",
	Storage:     ".",
	Directories: false,
	Port:        4100,
}

func init() {
	// viper.SetEnvPrefix("rake")
	viper.AutomaticEnv()
	viper.BindEnv("SECRET")
	viper.BindEnv("URI")
	viper.BindEnv("DATABASE")
	viper.BindEnv("STORAGE")
	viper.BindEnv("DIRECTORIES")
	viper.BindEnv("PORT")
	viper.SetConfigName(".rake")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
}
