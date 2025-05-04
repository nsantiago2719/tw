package main

import (
	"github.com/urfave/cli/v3"
)

var (
	initCommand = command{
		Name:    "init",
		Aliases: []string{"i"},
		Usage:   "initializes the current working directory as the parent directory",
		Action:  actionInit,
	}

	registerResource = command{
		Name:    "register",
		Aliases: []string{"r"},
		Usage:   "registers a resource to the config file",
		Action:  actionRegisterResource,
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

	resources = command{
		Name:    "list-resources",
		Aliases: []string{"lr"},
		Usage:   "list all resources",
		Action:  actionResources,
	}

	run = command{
		Name:  "run",
		Usage: "runs terraform apply against the resource values",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "dry-run",
				Usage: "sets the terraform dry-run flag",
			},
		},
		Action: actionRunTerraform,
	}

	plan = command{
		Name:  "plan",
		Usage: "run terraform plan against the resource values",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:      "resource-name",
				UsageText: "[resource-name]",
			},
		},
		Action: actionPlanTerraform,
	}
)
