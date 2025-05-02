package main

import (
	"os"
)

func main() {
	app := newApp()
	app.addCommand(initCommand)
	app.addCommand(registerResource)
	app.addCommand(resources)
	app.addCommand(run)
	app.addCommand(plan)
	app.run(os.Args)
}
