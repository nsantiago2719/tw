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
	Name        string
	Aliases     []string
	Usage       string
	Action      actionFunc
	Flags       []cli.Flag
	SubCommands []command
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

func (app *app) makeSubCommands(cmds []command) []*cli.Command {
	subCommands := []*cli.Command{}
	for _, cmd := range cmds {
		makeSubCommand := app.makeCommand(cmd)
		subCommands = append(subCommands, &makeSubCommand)
	}

	return subCommands
}

func (app *app) makeCommand(cmd command) cli.Command {
	return cli.Command{
		Name:        cmd.Name,
		Aliases:     cmd.Aliases,
		Usage:       cmd.Usage,
		Action:      makeAction(cmd.Action, app.configPath),
		Flags:       cmd.Flags,
		Subcommands: app.makeSubCommands(cmd.SubCommands),
	}
}

func (app *app) run(args []string) error {
	cliApp := &cli.App{
		Name:     app.name,
		Usage:    app.usage,
		Commands: app.commands,
	}

	return cliApp.Run(args)
}
