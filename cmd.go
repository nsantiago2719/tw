package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// This contains all command instruction that will be passed to terraform
type cmd struct {
	Command  string
	VarFiles []string
	DryRun   bool
}

func createCmd(command string, varFiles []string, dryRun bool) cmd {
	return cmd{
		Command:  command,
		VarFiles: varFiles,
		DryRun:   dryRun,
	}
}

func (cmd *cmd) addVarFile(file string) *cmd {
	cmd.VarFiles = append(cmd.VarFiles, file)

	return cmd
}

// TODO: run the said command using terraform
func (cmd *cmd) execCmd(ctx context.Context) error {
	var dryRun string
	varFiles := strings.Join(cmd.VarFiles, " ")
	if cmd.DryRun {
		dryRun = "-dry-run"
	}

	command := fmt.Sprintf("terraform %v %v %v", cmd.Command, varFiles, dryRun)

	exec.CommandContext(ctx, command)

	return nil
}
