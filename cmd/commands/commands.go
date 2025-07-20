package commands

import (
	"github.com/nsantiago2719/tw/internal/actions"
	"github.com/nsantiago2719/tw/internal/app"

	"github.com/urfave/cli/v3"
)

type commands []*app.Command

// AllCommands is the list of all available commands
var AllCommands = commands{
	&initCommand,
	&registerResource,
	&resources,
	&run,
	&plan,
}

var (
	initCommand = app.Command{
		Name:    "init",
		Aliases: []string{"i"},
		Usage:   "initializes the current working directory as the parent directory",
<<<<<<< HEAD:cmd/commands/commands.go
		Action:  actions.Init,
=======
		Action:  actionInit,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "opentofu",
				Usage: "Use opentofu instead of terraform.",
			},
		},
>>>>>>> 9e0188b (feat: added initial code for supporting opentofu instead of terraform):commands.go
	}

	registerResource = app.Command{
		Name:    "register",
		Aliases: []string{"r"},
		Usage:   "registers a resource to the config file",
		Action:  actions.RegisterResource,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Usage: "`resource-name` for the resource being added. (required: true)",
			},
			&cli.StringFlag{
				Name:  "path",
				Usage: "`PATH` where the resource is located. (required: true)",
			},
			&cli.StringSliceFlag{
				Name:  "var-files",
				Usage: "Load variable values from the given files.",
			},
		},
	}

	resources = app.Command{
		Name:    "list-resources",
		Aliases: []string{"lr"},
		Usage:   "list all resources",
		Action:  actions.Resources,
	}

	run = app.Command{
		Name:  "run",
		Usage: "runs terraform apply against the resource values",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:      "resource-name",
				UsageText: "[resource-name]",
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "dry-run",
				Usage: "sets the terraform dry-run flag",
			},
			&cli.BoolFlag{
				Name:  "auto-approve",
				Usage: "sets terraform auto-accept flag",
			},
		},
		Action: actions.RunTerraform,
	}

	plan = app.Command{
		Name:  "plan",
		Usage: "run terraform plan against the resource values",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:      "resource-name",
				UsageText: "[resource-name]",
			},
		},
		Action: actions.PlanTerraform,
	}

	destroy = command{
		Name:  "destroy",
		Usage: "run terraform destroy against the resource values",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:      "resource-name",
				UsageText: "[resource-name]",
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "auto-approve",
				Usage: "sets terraform auto-accept flag",
			},
			&cli.BoolFlag{
				Name:  "dry-run",
				Usage: "sets the terraform dry-run flag",
			},
		},
	}
)
