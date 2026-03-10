package main

import (
	"os"

	"github.com/cwbudde/go-citygml/cmd/citygml/cli"
)

func main() {
	err := cli.Execute()
	if err != nil {
		os.Exit(1)
	}
}
