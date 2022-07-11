package main

import (
	"os"

	"github.com/AppleGamer22/rake/cli/commands"
)

func main() {
	if err := commands.RootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
