package main

import (
	"context"

	"github.com/urfave/cli/v3"
)

type actionFunc func(context.Context, *cli.Command, string) error

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
		Name:     cmd.Name,
		Aliases:  cmd.Aliases,
		Usage:    cmd.Usage,
		Action:   makeAction(cmd.Action, app.configPath),
		Flags:    cmd.Flags,
		Commands: app.makeSubCommands(cmd.SubCommands),
	}
}

func (app *app) run(ctx context.Context, args []string) error {
	cliApp := &cli.Command{
		Name:     app.name,
		Usage:    app.usage,
		Commands: app.commands,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Load config file from `FILE`",
				DefaultText: "config.json",
			},
		},
	}

	return cliApp.Run(ctx, args)
}
