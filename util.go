package main

import (
	"context"
	"fmt"
	"log/slog"
)

func runInit(ctx context.Context, resourcePath string) error {
	argsInit := createArgs("init", resourcePath, []string{})
	execInit := createCmd("init", argsInit)

	execInitOutput, err := execInit.exec(ctx)
	if err != nil {
		return err
	}

	stdOutput(execInitOutput)

	return nil
}

func stdOutput(outLine <-chan stdOutLine) {
	for output := range outLine {
		var level string
		if output.Stream == "stderr" {
			level = fmt.Sprintf("level=%v", "ERROR")
		} else {
			level = fmt.Sprintf("level=%v", "INFO")
		}
		slog.Info(level, "msg", output.Msg)
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
