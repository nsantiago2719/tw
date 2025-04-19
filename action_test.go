package main

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cliApp = newApp()

func init() {
	cliApp.addCommand(initCommand)
	cliApp.addCommand(registerResource)
	cliApp.addCommand(resources)
}

func TestActionInit(t *testing.T) {
	// Run init(i) command to create the config file
	err := cliApp.run([]string{"tw", "i"})
	assert.NoError(t, err)

	// remove config json at the end of the test
	defer os.Remove(cliApp.configPath)

	// Check if the config file exist
	if _, err := os.Stat(cliApp.configPath); errors.Is(err, os.ErrNotExist) {
		assert.Nil(t, err.Error())
	}
}

func TestActionRegisterResource(t *testing.T) {
	// Create the config.json file
	cliApp.run([]string{"tw", "i"})

	// always remove the config.json
	defer os.Remove(cliApp.configPath)

	err := cliApp.run([]string{"tw", "r", "--name", "resource-name", "--path", "./resource-test"})
	assert.NoError(t, err)

	resources := []resource{}
	config, err := os.ReadFile("config.json")
	assert.NoError(t, err)
	err = json.Unmarshal(config, &resources)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(resources), "config.json should contain 1 resource")

	err = cliApp.run([]string{"tw", "r", "--name", "resource-name", "--path", "./resource-test"})
	assert.NoError(t, err)

	err = cliApp.run([]string{"tw", "r", "--name", "resource-name", "--path", "./resource-test", "--var-files", "./data/resource-a/values.tfvars", "--var-files", "./data/resource-a/values-b.tfvars"})
	assert.NoError(t, err)
}

func TestActionResources(t *testing.T) {
	// Create the config.json file
	cliApp.run([]string{"tw", "i"})

	mockResource := resource{
		Name: "resource-name",
		Path: "./resource-test",
	}
	mockResources := []resource{}

	mockResources = append(mockResources, mockResource)

	// always remove the config.json
	defer os.Remove(cliApp.configPath)

	err := cliApp.run([]string{"tw", "r", "--name", "resource-name", "--path", "./resource-test"})
	assert.NoError(t, err)

	resources := []resource{}
	config, err := os.ReadFile(cliApp.configPath)
	assert.NoError(t, err)

	err = json.Unmarshal(config, &resources)

	assert.NoError(t, err)

	err = cliApp.run([]string{"tw", "r", "--name", "resource-name", "--path", "./resource-test", "--var-files", "./data/resource-a/values.tfvars", "--var-files", "./data/resource-a/values-b.tfvars"})
	assert.NoError(t, err)

	err = cliApp.run([]string{"tw", "lr"})
	assert.Nil(t, err)
}

func TestActionRunTerraform(t *testing.T) {
	err := ""
	assert.Nil(t, err)
}
