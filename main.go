package main

import (
	"github.com/ctfang/command"
	"github.com/ctfang/goworker/console"
)

func main() {
	app := command.New()

	app.SetConfig("config.ini")
	app.IniConfig()

	AddCommands(&app)
	app.Run()
}

func AddCommands(app *command.Console) {
	app.AddCommand(&console.Register{})
	app.AddCommand(&console.Gateway{})
	app.AddCommand(&console.Worker{})
}
