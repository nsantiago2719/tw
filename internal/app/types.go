package app

import (
	"context"

	"github.com/urfave/cli/v3"
)

type app struct {
	name       string
	usage      string
	configPath string
	commands   []*cli.Command
}

type Command struct {
	Name        string
	Aliases     []string
	Usage       string
	Action      ActionFunc
	Arguments   []cli.Argument
	Flags       []cli.Flag
	SubCommands []Command
}
type Resource struct {
	Name     string   `json:"name"`
	Path     string   `json:"path"`
	VarFiles []string `json:"var-files"`
}
type StdOutLine struct {
	Stream string
	Msg    string
}

type ActionFunc func(context.Context, *cli.Command, string) error
