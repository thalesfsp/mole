package main

import (
	"os"

	"github.com/thalesfsp/mole/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
