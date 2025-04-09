package main

import (
	"errors"
	"os"
	"testing"

	"github.com/urfave/cli/v2"
)

func TestActionInit(t *testing.T) {
	t.Parallel()
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
	err := app.Run([]string{"tw", "i"})
	if err != nil {
		t.Fatalf("Test actionInit failed: %v", err)
	}

	if _, err := os.Stat("config.json"); errors.Is(err, os.ErrNotExist) {
		t.Errorf("Expected: nil, Got: %q", err)
	}
}
