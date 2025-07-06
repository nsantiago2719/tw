package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
)

// This contains all command instruction that will be passed to terraform
// Command is reffered to the terraform command instead of the command
// needed by the exec.CommandContext() function needs
// Args are the arguements passed to the command eg.  apply, plan
type cmd struct {
	Command string
	Args    []string
}

type stdOutLine struct {
	Stream string
	Msg    string
}

func initCmd(command string) cmd {
	return cmd{
		Command: command,
		Args:    []string{},
	}
}

func (cmd *cmd) createCmd(path string, varFiles ...string) error {
	chDir := fmt.Sprintf("-chdir=%v", path)
	// inject the chdir flag
	cmd.addArg(chDir)
	// inject the command eg. plan or apply
	cmd.addArg(cmd.Command)
	// inject no-color flag to remove ascii on the output
	cmd.addArg("-no-color")
	// inject var-files
	if len(varFiles) > 0 {
		for _, v := range varFiles {
			arg := fmt.Sprintf("-var-file=%v", v)
			cmd.addArg(arg)
		}
	}

	return nil
}

func (cmd *cmd) addArg(arg string) *cmd {
	cmd.Args = append(cmd.Args, arg)

	return cmd
}

func (cmd *cmd) exec(ctx context.Context) (<-chan stdOutLine, <-chan bool, chan<- string, error) {
	cmdCtx := exec.CommandContext(ctx, "terraform", cmd.Args...)

	stdoutPipe, err := cmdCtx.StdoutPipe()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create stdoutpipe: %w", err)
	}

	stderrPipe, err := cmdCtx.StderrPipe()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create stderrpipe: %w", err)
	}

	stdinPipe, err := cmdCtx.StdinPipe()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create stdinpipe: %w", err)
	}

	stdOutputChan := make(chan stdOutLine)
	stdinRequestChan := make(chan bool) // Signal when input is likely needed
	stdinInputChan := make(chan string) // Channel to receive user input

	var wg sync.WaitGroup

	// create readFromPipe func, gets the name and pipe
	readFromPipe := func(pipeName string, pipe io.ReadCloser) {
		defer wg.Done()
		scanner := bufio.NewScanner(pipe)
		for scanner.Scan() {
			line := scanner.Text()

			// Remove specific prompt as a workaround
			// since it is displayed after the input
			if contains(line, "Enter a value:") {
				line = ""
			}
			// do not send empty lines
			if line != "" {

				stdOutputChan <- stdOutLine{Stream: pipeName, Msg: line}

				// Check for common prompts that indicate the command is waiting for input
				// This is a heuristic approach - adjust these patterns based on your specific use case
				if containsInputPrompt(line) {
					select {
					case stdinRequestChan <- true:
						// Signal that input is needed
					default:
						// Channel is full or closed, do nothing
					}
				}
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("error reading %s: %v\n", pipeName, err)
		}
	}

	wg.Add(2)
	go readFromPipe("stdout", stdoutPipe)
	go readFromPipe("stderr", stderrPipe)

	// Handle stdin input from user
	go func() {
		for input := range stdinInputChan {
			_, err := io.WriteString(stdinPipe, input+"\n")
			if err != nil {
				fmt.Printf("error writing to stdin: %v\n", err)
				break
			}
		}
	}()

	err = cmdCtx.Start()
	if err != nil {
		close(stdOutputChan)
		close(stdinRequestChan)
		return nil, nil, nil, fmt.Errorf("failed to run command: %w", err)
	}
	go func() {
		err = cmdCtx.Wait()
		if err != nil {
			fmt.Printf("command finished with an error: %v\n", err)
		}
		wg.Wait()
		close(stdOutputChan)
		close(stdinRequestChan)
	}()

	return stdOutputChan, stdinRequestChan, stdinInputChan, nil
}

// containsInputPrompt checks if a line contains common patterns that indicate
// the command is waiting for user input
func containsInputPrompt(line string) bool {
	// Common patterns that might indicate a prompt for user input
	// Adjust these patterns based on the specific commands you're running
	inputPromptPatterns := []string{
		"Do you want to perform these actions?",
	}

	for _, pattern := range inputPromptPatterns {
		if contains(line, pattern) {
			return true
		}
	}
	return false
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	s, substr = strings.ToLower(s), strings.ToLower(substr)
	return strings.Contains(s, substr)
}
