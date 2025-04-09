package main

import (
	"fmt"
	"log/slog"

	"github.com/urfave/cli/v2"
)

func makeAction(f cli.ActionFunc) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		if err := f(ctx); err != nil {
			slog.Error("Error encountered: ", err)
		}
		return nil
	}
}

func actionInit(ctx *cli.Context) error {
	fmt.Println("initialize repository")

	return nil
}
