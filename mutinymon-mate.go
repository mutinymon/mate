package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = "mutinymon-mate"
	app.Version = Version
	app.Usage = ""
	app.Author = "Donovan Tengblad"
	app.Email = "purplefish32@gmail.com"
	app.Commands = Commands

	app.Run(os.Args)
}
