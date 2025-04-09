package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cli_app = newApp()

func TestActionInit(t *testing.T) {
	cli_app.addCommand(&initCommand)
	// Run init(i) command to create the config file
	err := cli_app.run([]string{"tw", "i"})
	if err != nil {
		t.Fatalf("Test actionInit failed: %v", err)
	}

	// remove config json at the end of the test
	defer os.Remove("config.json")

	// Check if the config file exist
	if _, err := os.Stat("config.json"); errors.Is(err, os.ErrNotExist) {
		assert.Nil(t, err.Error())
	}
}

func TestActionRegisterResource(t *testing.T) {
	cli_app.addCommand(&registerResource)
	// Create the config.json file
	cli_app.run([]string{"tw", "i"})

	// always remove the config.json
	defer os.Remove("config.json")

	err := cli_app.run([]string{"tw", "r", "--name", "resource-name", "--path", "./resource-test"})
	if err != nil {
		t.Errorf("Test actionRegisterResource failed: %v", err)
	}
}
