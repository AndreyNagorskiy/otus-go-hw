package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

const (
	SuccessCode = iota
	ErrorCode
	FileNotFoundCOde
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return ErrorCode
	}

	for key, envItem := range env {
		if envItem.NeedRemove {
			err := os.Unsetenv(key)
			if err != nil {
				return ErrorCode
			}
		}

		err := os.Setenv(key, envItem.Value)
		if err != nil {
			return ErrorCode
		}
	}

	path, err := exec.LookPath(cmd[0])
	if err != nil {
		fmt.Printf("Error: command %q not found\n", cmd[0])
		return ErrorCode
	}

	command := exec.Command(path, cmd[1:]...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err = command.Run()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}
		fmt.Printf("Command execution error %s: %v\n", cmd[0], err)
		return ErrorCode
	}

	return SuccessCode
}
