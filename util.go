package main

import (
	"context"
	"fmt"
	"log/slog"
)

func runInit(ctx context.Context, resourcePath string) error {
	execInit := initCmd("init")
	execInit.createCmd(resourcePath)

	execInitOutput, stdinRequestChan, stdinInputChan, err := execInit.exec(ctx)
	if err != nil {
		return err
	}

	// Use handleCommandIO to process both output and stdin requests
	handleCommandIO(execInitOutput, stdinRequestChan, stdinInputChan)

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

// Set a better way on getting the resouce details which allows to
// having the same resource name
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
