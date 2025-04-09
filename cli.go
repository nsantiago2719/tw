package main

import (
	"github.com/urfave/cli/v2"
)

type app struct {
	name     string
	usage    string
	commands []*cli.Command
}

func newApp() app {
	return app{
		name:  "tw",
		usage: "tw [commands]",
	}
}

func (app *app) addCommand(command *cli.Command) *app {
	app.commands = append(app.commands, command)
	return app
}

func (app *app) run(a []string) error {
	cli_app := &cli.App{
		Name:     app.name,
		Usage:    app.usage,
		Commands: app.commands,
	}

	return cli_app.Run(a)
}
