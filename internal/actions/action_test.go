package actions_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/nsantiago2719/tw/internal/app"
	"github.com/nsantiago2719/tw/internal/commands"
	"github.com/stretchr/testify/assert"
)

var (
	cliApp = app.NewApp()
	ctx    = context.Background()
)

func init() {
	cliApp.AddCommand(commands.InitCommand)
	cliApp.AddCommand(commands.RegisterResource)
	cliApp.AddCommand(commands.Resources)
	cliApp.AddCommand(commands.Run)
}

func TestActionInit(t *testing.T) {
	// Run init(i) command to create the config file
	err := cliApp.Run(ctx, []string{"tw", "i"})
	assert.NoError(t, err)

	// always remove the config.json
	defer os.Remove(cliApp.ConfigPath())

	// Check if the config file exist
	if _, err := os.Stat(cliApp.ConfigPath()); errors.Is(err, os.ErrNotExist) {
		assert.Nil(t, err.Error())
	}
}

// TODO: TestActionPlanTerraform
// TODO: TestActionRunTerraform

func TestActionRegisterResource(t *testing.T) {
	// Create the config.json file
	cliApp.Run(ctx, []string{"tw", "i"})

	// always remove the config.json
	defer os.Remove(cliApp.ConfigPath())

	err := cliApp.Run(ctx, []string{"tw", "r", "--name", "resource-name", "--path", "./resource-test"})
	assert.NoError(t, err)

	resources := []app.Resource{}
	config, err := os.ReadFile("config.json")
	assert.NoError(t, err)
	err = json.Unmarshal(config, &resources)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(resources), "config.json should contain 1 resource")

	err = cliApp.Run(ctx, []string{"tw", "r", "--name", "resource-name", "--path", "./resource-test"})
	assert.NoError(t, err)

	err = cliApp.Run(ctx, []string{"tw", "r", "--name", "resource-name", "--path", "./resource-test", "--var-files", "../../test/resource-a/first.tfvars", "--var-files", "../../test/resource-a/second.tfvars"})
	assert.NoError(t, err)
}

func TestActionResources(t *testing.T) {
	// Create the config.json file
	cliApp.Run(ctx, []string{"tw", "i"})

	mockResource := app.Resource{
		Name: "resource-name",
		Path: "./resource-test",
	}
	mockResources := []app.Resource{}

	mockResources = append(mockResources, mockResource)

	// always remove the config.json
	defer os.Remove(cliApp.ConfigPath())

	err := cliApp.Run(ctx, []string{"tw", "r", "--name", "resource-name", "--path", "../../test"})
	assert.NoError(t, err)

	resources := []app.Resource{}
	config, err := os.ReadFile(cliApp.ConfigPath())
	assert.NoError(t, err)

	err = json.Unmarshal(config, &resources)

	assert.NoError(t, err)

	err = cliApp.Run(ctx, []string{"tw", "r", "--name", "resource-name", "--path", "../../test/resource-a", "--var-files", "../../test/resource-a/first.tfvars", "--var-files", "../../test/resource-a/second.tfvars"})
	assert.NoError(t, err)

	err = cliApp.Run(ctx, []string{"tw", "lr"})
	assert.Nil(t, err)
}

func TestActionRunTerraform(t *testing.T) {
	cliApp.Run(ctx, []string{"tw", "i"})
	err := cliApp.Run(ctx, []string{"tw", "r", "--name", "resource-name", "--path", "../../test/resource-a", "--var-files", "first.tfvars", "--var-files", "second.tfvars"})
	assert.NoError(t, err)

	// always remove the config.json
	defer os.Remove(cliApp.ConfigPath())

	err = cliApp.Run(ctx, []string{"tw", "run", "resource-name", "-auto-approve"})
	assert.NoError(t, err)
}
