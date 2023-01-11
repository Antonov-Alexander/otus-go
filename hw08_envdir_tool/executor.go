package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for key, value := range env {
		var err error
		if !value.NeedRemove {
			err = os.Setenv(key, value.Value)
		} else {
			err = os.Unsetenv(key)
		}

		if err != nil {
			return 1
		}
	}

	executor := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	executor.Env = os.Environ()
	executor.Stdin = os.Stdin
	executor.Stdout = os.Stdout
	executor.Stderr = os.Stderr
	if err := executor.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}

		return 1
	}

	return 0
}
