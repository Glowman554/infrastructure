package utils

import (
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

func createCommand(command string) *exec.Cmd {
	tmp := strings.Split(command, " ")
	tmp = Filter(tmp, func(s string) bool {
		return s != ""
	})
	slog.Debug("executing command", "command", tmp)
	return exec.Command(tmp[0], tmp[1:]...)
}

func Execute(command string, directory string, show bool) error {
	cmd := createCommand(command)
	cmd.Dir = directory

	out := strings.Builder{}

	cmd.Stdout = &out
	cmd.Stderr = &out
	if show {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	err := cmd.Run()
	if err != nil {
		slog.Error("command failed", "command", command, "error", err, "output", out.String())
		return err
	}

	return nil
}
