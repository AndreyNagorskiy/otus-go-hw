package main

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	ErrorCode = iota
	SuccessCode
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
		fmt.Printf("Ошибка: команда %q не найдена\n", cmd[0])
		return ErrorCode
	}

	command := exec.Command(path, cmd[1:]...)

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err = command.Run()
	if err != nil {
		fmt.Printf("Command execution error %s: %v\n", cmd[0], err)
		return ErrorCode
	}

	return SuccessCode
}
