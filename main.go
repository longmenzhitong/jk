package main

import (
	"fmt"
	"jk/cmd"
	jk "jk/src"
	"os"
)

func main() {
	jk.Config.Init()

	if err := cmd.RootCmd.Execute(); err != nil {
		if err != cmd.SilentErr {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
