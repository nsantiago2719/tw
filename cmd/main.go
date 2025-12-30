package main

import (
	"context"
	"os"

	"github.com/nsantiago2719/tw/cmd/commands"
	"github.com/nsantiago2719/tw/internal/app"
)

func main() {
	ctx := context.Background()

	cli := app.NewApp()

	// add commands iteratively
	for _, cmd := range commands.AllCommands {
		cli.AddCommand(*cmd)
	}
	cli.Run(ctx, os.Args)
}
