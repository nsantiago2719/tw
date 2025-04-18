package main

import (
	"github.com/urfave/cli/v2"
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
				Value: "resoure-name",
				Usage: "the name of the resource to be added",
			},
			&cli.StringFlag{
				Name:  "path",
				Value: "resource-path",
				Usage: "directory path of the resource",
			},
			&cli.StringSliceFlag{
				Name:  "var-files",
				Usage: "var file path for the resource if needed",
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
		Name:   "run",
		Usage:  "runs terraform against the resource values",
		Action: actionRunTerraform,
	}
)
