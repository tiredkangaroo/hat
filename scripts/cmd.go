package main

import (
	"log/slog"
	"os"
	"os/exec"

	"github.com/fatih/color"
)

func cmderror(s string) error {
	color.Blue("shell: %s\n", s)
	cmd := exec.Command("bash", "-c", s)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func cmd(s string) {
	if err := cmderror(s); err != nil {
		slog.Error("command failed", "error", err)
		os.Exit(1)
	}
}
