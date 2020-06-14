package main

import (
	"os"
	"path/filepath"

	"github.com/UiP9AV6Y/changelog-assembler/pkg/cmd"
)

func main() {
	app := filepath.Base(os.Args[0])
	root := cmd.New(app)

	if err := root.Command().Execute(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
