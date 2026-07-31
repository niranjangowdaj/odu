package main

import (
	"fmt"
	"os"

	"github.com/odu-cli/odu/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
