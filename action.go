package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v2"
)

type resource struct {
	Name     string   `json:"name"`
	Path     string   `json:"path"`
	VarFiles []string `json:"var-files"`
}

// makeAction is a wrapper for injecting generic code for all actions
// eg. logging
func makeAction(f cli.ActionFunc) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		if err := f(ctx); err != nil {
			slog.Error("level=error", "msg", err)
			return err
		}
		return nil
	}
}

func actionRegisterResource(ctx *cli.Context) error {
	if ctx.String("name") == "" {
		return errors.New("Name must not be empty")
	}

	if ctx.String("path") == "" {
		return errors.New("Path must not be empty")
	}

	rs := resource{
		Name:     ctx.String("name"),
		Path:     ctx.String("path"),
		VarFiles: ctx.StringSlice("var-files"),
	}

	file, err := os.OpenFile("config.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil
	}

	defer file.Close()

	resources := []resource{}

	config, err := os.ReadFile("config.json")
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

func actionResources(ctx *cli.Context) error {
	config, err := os.ReadFile("config.json")
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
func actionInit(ctx *cli.Context) error {
	if _, err := os.Stat("config.json"); errors.Is(err, os.ErrNotExist) {
		err := os.WriteFile("config.json", []byte("[]"), 0755)
		if err != nil {
			return err
		}
	}
	fmt.Println("tw initialized: created config.json file")
	return nil
}
