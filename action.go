package main

import (
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"
)

func makeAction(f cli.ActionFunc) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		if err := f(ctx); err != nil {
			slog.Error("level=error", "msg", err)
			return err
		}
		return nil
	}
}

func actionInit(ctx *cli.Context) error {
	err := os.WriteFile("config.json", []byte(""), 0755)
	if err != nil {
		return err
	}
	return nil
}
