package main

import (
	"os"

	"github.com/nyarly/git-along/along"
)

func main() {
	if err := along.Execute(); err != nil {
		os.Exit(1)
	}
}
