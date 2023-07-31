package main

import (
	"jk/cmd"
	jk "jk/src"
)

func main() {
	jk.Config.Init()

	cmd.Execute()
}
