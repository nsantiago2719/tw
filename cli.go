package main

import (
	"github.com/urfave/cli/v2"
)

type actionFunc func(*cli.Context, string) error

type app struct {
	name       string
	usage      string
	configPath string
	commands   []*cli.Command
}

type command struct {
	Name    string
	Aliases []string
	Usage   string
	Action  actionFunc
	Flags   []cli.Flag
}

func newApp() app {
	return app{
		name:       "tw",
		configPath: "config.json",
		usage:      "tw [commands]",
	}
}

func (app *app) addCommand(cmd command) *app {
	makeCmd := app.makeCommand(cmd)
	app.commands = append(app.commands, &makeCmd)
	return app
}

func (app *app) makeCommand(command command) cli.Command {
	return cli.Command{
		Name:    command.Name,
		Aliases: command.Aliases,
		Usage:   command.Usage,
		Action:  makeAction(command.Action, app.configPath),
		Flags:   command.Flags,
	}
}

func (app *app) run(a []string) error {
	cli_app := &cli.App{
		Name:     app.name,
		Usage:    app.usage,
		Commands: app.commands,
	}

	return cli_app.Run(a)
}
