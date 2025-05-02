package main

import (
	"context"
	"os"
)

func main() {
	ctx := context.Background()

	app := newApp()
	app.addCommand(initCommand)
	app.addCommand(registerResource)
	app.addCommand(resources)
	app.addCommand(run)
	app.addCommand(plan)
	app.run(ctx, os.Args)
}
