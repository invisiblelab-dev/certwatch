package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/invisiblelab-dev/certwatch/commands"
)

func main() {
	c := commands.Parse()
	if err := c.Execute(); err != nil {
		if !errors.Is(err, commands.ErrSilent) {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
