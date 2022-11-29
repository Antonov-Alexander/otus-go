package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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

	stdout, stderr := new(strings.Builder), new(strings.Builder)
	executor := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	executor.Env = os.Environ()
	executor.Stdout = stdout
	executor.Stderr = stderr
	if err := executor.Run(); err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
		return 1
	}

	if _, err := fmt.Fprint(os.Stdout, stdout.String()); err != nil {
		return 1
	}

	return 0
}
