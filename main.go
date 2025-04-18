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
	app.run(os.Args)
}
