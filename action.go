package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"
)

type Resource struct {
	Name string `json:name`
	Path string `json:path`
}

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

func actionRegisterResource(ctx *cli.Context) error {
	return os.ErrNotExist
}

// actionInit create a config.json file if the file does not exist
// else it would do nothing
func actionInit(ctx *cli.Context) error {
	if _, err := os.Stat("config.json"); errors.Is(err, os.ErrNotExist) {
		err := os.WriteFile("config.json", []byte("[]"), 0755)
		if err != nil {
			return err
		}
	}
	fmt.Println("tw initialized: created config.json file")
	return nil
}
