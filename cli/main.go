package main

import (
	"log"
	"os"

	"github.com/AppleGamer22/rake/cli/commands"
	"github.com/spf13/viper"
)

func init() {
	viper.SetEnvPrefix("rake")
	viper.AutomaticEnv()
	viper.SetConfigName(".rake")
	viper.SetConfigType("yaml")

	directory, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	viper.AddConfigPath(directory)

	// session, fbsr & app
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	if err := commands.RootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
