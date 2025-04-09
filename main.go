package main

import (
	"os"
)

func main() {
	app := newApp()
	app.addCommand(&initCommand)
	app.addCommand(&registerResource)
	app.run(os.Args)
}
