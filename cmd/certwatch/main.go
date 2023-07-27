package main

import (
	"github.com/invisiblelab-dev/certwatch/commands"
)

func main() {
	c := commands.Parse()
	_ = c.Execute()
}
