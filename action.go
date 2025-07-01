package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v3"
)

type resource struct {
	Name     string   `json:"name"`
	Path     string   `json:"path"`
	VarFiles []string `json:"var-files"`
}

// makeAction is a wrapper for injecting generic code for all actions
// eg. logging
func makeAction(f actionFunc, cfg string) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		// sets the to default cfg if config flag is not passed
		var cfgPath string
		if cmd.String("config") == "" {
			cfgPath = cfg
		} else {
			cfgPath = cmd.String("config")
		}
		if err := f(ctx, cmd, cfgPath); err != nil {
			slog.Error("level=error", "msg", err)
			return err
		}
		return nil
	}
}

// printOutput is a local helper function to print command output with color formatting
func printOutput(output <-chan stdOutLine) {
	for line := range output {
		if line.Stream == "stderr" {
			fmt.Printf("\033[31m%s\033[0m\n", line.Msg)
		} else {
			fmt.Println(line.Msg)
		}
	}
}

// handleStdin manages user input when a command requires it
func handleStdin(stdinRequestChan <-chan bool, stdinInputChan chan<- string) {
	for range stdinRequestChan {
		// Command is waiting for input
		fmt.Print("\033[33mInput required: \033[0m") // Yellow prompt

		// Read user input
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			userInput := scanner.Text()
			// Send the input to the command
			stdinInputChan <- userInput
		} else {
			// Handle scanner error
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", scanner.Err())
			break
		}
	}
}

func handleCommandIO(output <-chan stdOutLine, stdinRequestChan <-chan bool, stdinInputChan chan<- string) {
	// Start a goroutine to handle stdin requests
	go handleStdin(stdinRequestChan, stdinInputChan)

	// Process output in the current goroutine
	printOutput(output)
}

// TODO: run terraform apply along with the var files passed if exist
func actionRunTerraform(ctx context.Context, cmd *cli.Command, cfg string) error {
	resourceName := cmd.StringArg("resource-name")
	if resourceName == "" {
		return fmt.Errorf("resource-name cannot be empty")
	}

	var resources []resource
	config, err := os.ReadFile(cfg)
	if err != nil {
		return err
	}

	err = json.Unmarshal(config, &resources)
	if err != nil {
		return err
	}

	resourcePath, args := getDetails(resourceName, resources)
	if resourcePath == "" {
		return fmt.Errorf("the resource registered has an empty path")
	}

	if err := runInit(ctx, resourcePath); err != nil {
		return err
	}

	execApply := initCmd("apply")
	err = execApply.createCmd(resourcePath, args...)
	if err != nil {
		return err
	}

	if cmd.Bool("auto-approve") {
		execApply.addArg("-auto-approve")
	}

	if cmd.Bool("dry-run") {
		execApply.addArg("-dry-run")
	}

	// Get output channel and stdin channels from command execution
	// The stdinRequestChan and stdinInputChan are used by handleCommandIO to manage interactive input
	execApplyOutput, stdinRequestChan, stdinInputChan, err := execApply.exec(ctx)
	if err != nil {
		return err
	}

	// Handle both output and potential input requests
	// This function uses all three channels to manage command IO
	handleCommandIO(execApplyOutput, stdinRequestChan, stdinInputChan)

	return nil
}

func actionPlanTerraform(ctx context.Context, cmd *cli.Command, cfg string) error {
	resourceName := cmd.StringArg("resource-name")

	if resourceName == "" {
		return fmt.Errorf("resouce-name cannot be empty")
	}

	var resources []resource
	config, err := os.ReadFile(cfg)
	if err != nil {
		return err
	}

	err = json.Unmarshal(config, &resources)
	if err != nil {
		return err
	}

	resourcePath, varFiles := getDetails(resourceName, resources)
	if resourcePath == "" {
		return fmt.Errorf("the resource registered has an empty path")
	}

	if err := runInit(ctx, resourcePath); err != nil {
		return err
	}

	execPlan := initCmd("plan")
	err = execPlan.createCmd(resourcePath, varFiles...)
	if err != nil {
		return err
	}

	// Get output channel and stdin channels from command execution
	// The stdinRequestChan and stdinInputChan are used by handleCommandIO to manage interactive input
	execPlanOutput, stdinRequestChan, stdinInputChan, err := execPlan.exec(ctx)
	if err != nil {
		return err
	}

	// Handle both output and potential input requests
	// This function uses all three channels to manage command IO
	handleCommandIO(execPlanOutput, stdinRequestChan, stdinInputChan)

	return nil
}

// TODO: set path as the current path if path flag is `.`
func actionRegisterResource(_ context.Context, cmd *cli.Command, cfg string) error {
	if cmd.String("name") == "" {
		return errors.New("Name must not be empty")
	}

	if cmd.String("path") == "" {
		return errors.New("Path must not be empty")
	}

	rs := resource{
		Name:     cmd.String("name"),
		Path:     cmd.String("path"),
		VarFiles: cmd.StringSlice("var-files"),
	}

	file, err := os.OpenFile(cfg, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil
	}

	defer file.Close()

	resources := []resource{}

	config, err := os.ReadFile(cfg)
	if err != nil {
		return err
	}

	err = json.Unmarshal(config, &resources)
	if err != nil {
		return nil
	}

	for resource := range resources {
		nameIndex := resource + 1
		baseName := cmd.String("name")
		if resources[resource].Name == rs.Name {
			rs.Name = fmt.Sprintf("%s-%d", baseName, nameIndex)
		}
	}

	resources = append(resources, rs)

	marshalResources, err := json.MarshalIndent(resources, "", "  ")
	if err != nil {
		return err
	}

	err = file.Truncate(0)
	if err != nil {
		return err
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = file.Write(marshalResources)
	if err != nil {
		return err
	}

	return nil
}

func actionResources(_ context.Context, _ *cli.Command, cfg string) error {
	config, err := os.ReadFile(cfg)
	resources := []resource{}
	if err != nil {
		return err
	}

	err = json.Unmarshal(config, &resources)

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Name", "Path", "Var Files")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, resource := range resources {
		tbl.AddRow(resource.Name, resource.Path, resource.VarFiles)
	}

	tbl.Print()
	return nil
}

// actionInit create a config.json file if the file does not exist
// else it would do nothing
func actionInit(_ context.Context, _ *cli.Command, cfg string) error {
	if _, err := os.Stat(cfg); errors.Is(err, os.ErrNotExist) {
		err := os.WriteFile(cfg, []byte("[]"), 0755)
		if err != nil {
			return err
		}
	}
	fmt.Println("tw initialized: created config.json file")
	return nil
}
