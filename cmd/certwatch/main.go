package main

import (
	"github.com/invisiblelab-dev/certwatch/internal/commands"
)

func main() {
	c := commands.Parse()
	_ = c.Execute()
}