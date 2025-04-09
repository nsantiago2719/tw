package main

import (
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "tw",
		Usage: "tw [commands]",
		Commands: []*cli.Command{
			{
				Name:    "init",
				Aliases: []string{"i"},
				Usage:   "initializes the current working diretory as the parent directory",
				Action:  makeAction(actionInit),
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("Error encountered:", err)
	}
}
