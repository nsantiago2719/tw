package main

import (
	"context"
	"os"

	"github.com/nsantiago2719/tw/internal/app"
	"github.com/nsantiago2719/tw/internal/commands"
)

func main() {
	ctx := context.Background()

	cli := app.NewApp()
	cli.AddCommand(commands.InitCommand)
	cli.AddCommand(commands.RegisterResource)
	cli.AddCommand(commands.Resources)
	cli.AddCommand(commands.Run)
	cli.AddCommand(commands.Plan)
	cli.Run(ctx, os.Args)
}
