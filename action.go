package main

import (
	"errors"
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"
)

// makeAction is a wrapper for injecting generic code for all actions
// eg. logging
func makeAction(f cli.ActionFunc) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		if err := f(ctx); err != nil {
			slog.Error("level=error", "msg", err)
			return err
		}
		return nil
	}
}

// actionInit create a config.json file if the file does not exist
// else it would do nothing
func actionInit(ctx *cli.Context) error {
	if _, err := os.Stat("config.json"); errors.Is(err, os.ErrNotExist) {
		err := os.WriteFile("config.json", []byte(""), 0755)
		if err != nil {
			return err
		}
	}

	return nil
}
