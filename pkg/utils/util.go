package utils

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/nsantiago2719/tw/internal/app"
	"github.com/nsantiago2719/tw/internal/terraform"
)

// RunAction executes the provided action function with the given context, command, and resource path.
func RunInit(ctx context.Context, resourcePath string) error {
	execInit := terraform.InitCmd("init")
	execInit.CreateCmd(resourcePath)

	execInitOutput, stdinRequestChan, stdinInputChan, err := execInit.Exec(ctx)
	if err != nil {
		return err
	}

	// Use handleCommandIO to process both output and stdin requests
	HandleCommandIO(execInitOutput, stdinRequestChan, stdinInputChan)

	return nil
}

// StdOutput processes standard output lines from a channel and logs them appropriately.
func StdOutput(outLine <-chan app.StdOutLine) {
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

// GetDetails retrieves the path and variable files for a given resource name from a list of resources.
func GetDetails(name string, resources []app.Resource) (string, []string) {
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

// printOutput is a local helper function to print command output with color formatting
func PrintOutput(output <-chan app.StdOutLine) {
	for line := range output {
		if line.Stream == "stderr" {
			fmt.Printf("\033[31m%s\033[0m\n", line.Msg)
		} else {
			fmt.Println(line.Msg)
		}
	}
}

// handleStdin manages user input when a command requires it
func HandleStdin(stdinRequestChan <-chan bool, stdinInputChan chan<- string) {
	for range stdinRequestChan {
		// Command is waiting for input
		fmt.Print("\n\033[33mInput required: \033[0m") // Yellow prompt

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

func HandleCommandIO(output <-chan app.StdOutLine, stdinRequestChan <-chan bool, stdinInputChan chan<- string) {
	// Start a goroutine to handle stdin requests
	go HandleStdin(stdinRequestChan, stdinInputChan)

	// Process output in the current goroutine
	PrintOutput(output)
}
