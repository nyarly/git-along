package main

import (
	"os"

	"git-along/along"
)

func main() {
	if err := along.Execute(); err != nil {
		os.Exit(1)
	}
}
