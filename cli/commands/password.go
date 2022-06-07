package commands

import (
	"fmt"
	"syscall"

	"golang.org/x/term"
)

func readPassword(password *string) error {
	bytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}
	fmt.Println()
	*password = string(bytes)
	return nil
}
