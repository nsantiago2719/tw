package commands

import (
	"github.com/nsantiago2719/tw/internal/actions"
	"github.com/nsantiago2719/tw/internal/app"
	"github.com/urfave/cli/v3"
)

var (
	InitCommand = app.Command{
		Name:    "init",
		Aliases: []string{"i"},
		Usage:   "initializes the current working directory as the parent directory",
		Action:  actions.Init,
	}

	RegisterResource = app.Command{
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

	Resources = app.Command{
		Name:    "list-resources",
		Aliases: []string{"lr"},
		Usage:   "list all resources",
		Action:  actions.Resources,
	}

	Run = app.Command{
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

	Plan = app.Command{
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
)
