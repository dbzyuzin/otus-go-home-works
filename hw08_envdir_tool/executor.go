package main

import (
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	command.Env = mergeEnvironments(os.Environ(), env)

	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			return exit.ExitCode()
		}
	}
	return 0
}

func mergeEnvironments(osenv []string, env Environment) []string {
	newEnv := make([]string, 0, len(osenv))

	for _, val := range osenv {
		name := strings.Split(val, "=")[0]
		nv, ok := env[name]
		if !ok {
			newEnv = append(newEnv, val)
			continue
		}
		delete(env, name)
		if nv == "" {
			continue
		}
		newEnv = append(newEnv, name+"="+nv)
	}

	for name, value := range env {
		newEnv = append(newEnv, name+"="+value)
	}
	return newEnv
}
