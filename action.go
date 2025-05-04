package main

import (
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

func getDetails(name string, resources []resource) (string, []string) {
	var path string
	var varFiles []string
	for _, v := range resources {
		if v.Name == name {
			path = v.Path
			varFiles = v.VarFiles
			// Stop looking for the resource
			// This will only allow us to get the first resource that
			// matches the resource name passed.
			break
		}
	}
	return path, varFiles
}

func createArgs(cmd string, path string, varFiles []string) []string {
	arg := []string{}
	chDir := fmt.Sprintf("-chdir=%v", path)
	arg = append(arg, chDir)
	// Append cmd after chdir since -chdir should be declared first
	arg = append(arg, cmd)
	arg = append(arg, "-no-color")
	if len(varFiles) > 0 {
		for _, v := range varFiles {
			text := fmt.Sprintf("-var-file=%v", v)
			arg = append(arg, text)
		}
	}

	return arg
}

// TODO: run terraform apply along with the var files passed if exist
func actionRunTerraform(ctx context.Context, cmd *cli.Command, cfg string) error {
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
		return fmt.Errorf("there is no resource registered with the name or the path is empty")
	}

	args := createArgs("plan", resourcePath, varFiles)
	execCommand := createCmd("plan", args)

	execCmd, err := execCommand.exec(ctx)
	if err != nil {
		return err
	}

	for output := range execCmd {
		var level string
		if output.Stream == "stderr" {
			level = fmt.Sprintf("level=%v", "ERROR")
		} else {
			level = fmt.Sprintf("level=%v", "INFO")
		}
		slog.Info(level, "msg", output.Msg)
	}

	return err
}

func actionRegisterResource(ctx context.Context, cmd *cli.Command, cfg string) error {
	if cmd.String("name") == "" {
		return errors.New("Name must not be empty")
	}

	if cmd.String("path") == "" {
		return errors.New("Path must not be empty")
	}

	fmt.Println(cmd.StringSlice("var-files"))
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

func actionResources(ctx context.Context, cmd *cli.Command, cfg string) error {
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
func actionInit(ctx context.Context, cmd *cli.Command, cfg string) error {
	if _, err := os.Stat(cfg); errors.Is(err, os.ErrNotExist) {
		err := os.WriteFile(cfg, []byte("[]"), 0755)
		if err != nil {
			return err
		}
	}
	fmt.Println("tw initialized: created config.json file")
	return nil
}
