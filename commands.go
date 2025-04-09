package main

import (
	"github.com/urfave/cli/v2"
)

var (
	initCommand = cli.Command{
		Name:    "init",
		Aliases: []string{"i"},
		Usage:   "initializes the current working directory as the parent directory",
		Action:  makeAction(actionInit),
	}

	registerResource = cli.Command{
		Name:    "register",
		Aliases: []string{"r"},
		Usage:   "registers a resource to the config file",
		Action:  makeAction(actionRegisterResource),
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
		},
	}
)
