package main

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cli_app = newApp()

func init() {
	cli_app.addCommand(&initCommand)
	cli_app.addCommand(&registerResource)
	cli_app.addCommand(&resources)
}

func TestActionInit(t *testing.T) {
	// Run init(i) command to create the config file
	err := cli_app.run([]string{"tw", "i"})
	assert.NoError(t, err)

	// remove config json at the end of the test
	defer os.Remove("config.json")

	// Check if the config file exist
	if _, err := os.Stat("config.json"); errors.Is(err, os.ErrNotExist) {
		assert.Nil(t, err.Error())
	}
}

func TestActionRegisterResource(t *testing.T) {
	// Create the config.json file
	cli_app.run([]string{"tw", "i"})

	// always remove the config.json
	defer os.Remove("config.json")

	err := cli_app.run([]string{"tw", "r", "--name", "resource-name", "--path", "./resource-test"})
	assert.NoError(t, err)

	resources := []resource{}
	config, err := os.ReadFile("config.json")
	assert.NoError(t, err)
	err = json.Unmarshal(config, &resources)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(resources), "config.json should contain 1 resource")

	err = cli_app.run([]string{"tw", "r", "--name", "resource-name", "--path", "./resource-test"})
	assert.NoError(t, err)
}

func TestActionResources(t *testing.T) {
	// Create the config.json file
	cli_app.run([]string{"tw", "i"})

	mockResource := resource{
		Name: "resource-name",
		Path: "./resource-test",
	}
	mockResources := []resource{}

	mockResources = append(mockResources, mockResource)

	// always remove the config.json
	defer os.Remove("config.json")

	err := cli_app.run([]string{"tw", "r", "--name", "resource-name", "--path", "./resource-test"})
	assert.NoError(t, err)

	resources := []resource{}
	config, err := os.ReadFile("config.json")
	assert.NoError(t, err)

	err = json.Unmarshal(config, &resources)

	assert.NoError(t, err)

	err = cli_app.run([]string{"tw", "lr"})
	assert.Nil(t, err)
}
