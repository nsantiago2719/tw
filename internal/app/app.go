package app

import (
	"context"
	"log/slog"

	"github.com/urfave/cli/v3"
)

func NewApp() app {
	return app{
		name:       "tw",
		configPath: "config.json",
		usage:      "tw [commands]",
	}
}

func (app *app) ConfigPath() string {
	return app.configPath
}

func (app *app) AddCommand(cmd Command) *app {
	makeCmd := app.makeCommand(cmd)
	app.commands = append(app.commands, &makeCmd)
	return app
}

func (app *app) makeSubCommands(cmds []Command) []*cli.Command {
	subCommands := []*cli.Command{}
	for _, cmd := range cmds {
		makeSubCommand := app.makeCommand(cmd)
		subCommands = append(subCommands, &makeSubCommand)
	}

	return subCommands
}

// makeAction is a wrapper for injecting generic code for all actions
// eg. logging
func makeAction(f ActionFunc, cfg string) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		// sets the to default cfg if config flag is not passed
		var cfgPath string
		if cmd.String("config") == "" {
			cfgPath = cfg
		} else {
			cfgPath = cmd.String("config")
		}
		if err := f(ctx, cmd, cfgPath); err != nil {
			slog.Error("level=error", "msg", err)
			return err
		}
		return nil
	}
}

func (app *app) makeCommand(cmd Command) cli.Command {
	return cli.Command{
		Name:      cmd.Name,
		Aliases:   cmd.Aliases,
		Usage:     cmd.Usage,
		Action:    makeAction(cmd.Action, app.configPath),
		Flags:     cmd.Flags,
		Arguments: cmd.Arguments,
		Commands:  app.makeSubCommands(cmd.SubCommands),
	}
}

func (app *app) Run(ctx context.Context, args []string) error {
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
