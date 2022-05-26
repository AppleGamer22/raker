package main

import (
	"os"

	"github.com/AppleGamer22/rake/cli/cmd"
)

func main() {
	if err := cmd.RootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
