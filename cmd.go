package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
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

func createCmd(command string, args []string) cmd {
	return cmd{
		Command: command,
		Args:    args,
	}
}

func (cmd *cmd) addArg(arg string) *cmd {
	cmd.Args = append(cmd.Args, arg)

	return cmd
}

// TODO: run the said command using terraform
// TODO: get the values from the config file
func (cmd *cmd) exec(ctx context.Context) (<-chan stdOutLine, error) {
	cmdCtx := exec.CommandContext(ctx, "terraform", cmd.Args...)

	stdoutPipe, err := cmdCtx.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdoutpipe: %w", err)
	}

	stderrPipe, err := cmdCtx.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderrpipe: %w", err)
	}

	stdOutputChan := make(chan stdOutLine)
	var wg sync.WaitGroup

	readFromPipe := func(pipeName string, pipe io.ReadCloser) {
		defer wg.Done()
		scanner := bufio.NewScanner(pipe)
		for scanner.Scan() {
			line := scanner.Text()
			if line != "" {
				stdOutputChan <- stdOutLine{Stream: pipeName, Msg: line}
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("error reading %s: %v\n", pipeName, err)
		}
	}

	wg.Add(1)
	go readFromPipe("stdout", stdoutPipe)

	wg.Add(1)
	go readFromPipe("stderr", stderrPipe)

	err = cmdCtx.Start()
	if err != nil {
		close(stdOutputChan)
		return nil, fmt.Errorf("failed to run command: %w", err)
	}

	go func() {
		err = cmdCtx.Wait()
		if err != nil {
			fmt.Errorf("command finished with an error: %w", err)
		}
		wg.Wait()
		defer close(stdOutputChan)
	}()

	return stdOutputChan, nil
}
