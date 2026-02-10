package main

import (
	"os"

	"github.com/ashhatz/launch-pad/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
