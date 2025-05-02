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
		cfgPath := ""
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

// TODO: run terraform apply along with the var files passed if exist
func actionRunTerraform(ctx context.Context, cmd *cli.Command, cfg string) error {
	return nil
}

func actionPlanTerraform(ctx context.Context, cmd *cli.Command, cfg string) error {
	// TODO: fetch varfiles from config
	// TODO: add dry run flag
	varFiles := []string{"-var-file=./hello/path/var.tfvars", "-var-file=./hello/path/extend.tfvars"}
	execCommand := createCmd("plan", varFiles, true)

	context := context.Background()

	return execCommand.execCmd(context)
}

func actionRegisterResource(ctx context.Context, cmd *cli.Command, cfg string) error {
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
