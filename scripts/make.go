package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type Command int

const (
	CommandNone Command = iota
	CommandRunDebug
	CommandRun
)

var currentCommand Command
var args []string

func init() {
	initLogger()
	if runtime.GOOS != "darwin" && runtime.GOOS != "linux" {
		fmt.Println("unsupported OS")
	}

	flag.Parse()
	args = flag.Args()

	if len(args) == 0 {
		return
	}

	switch args[0] {
	case "run":
		currentCommand = CommandRun
	case "run-debug":
		currentCommand = CommandRunDebug
	}
}

func main() {
	var err error
	switch currentCommand {
	case CommandRunDebug:
		err = runDebug()
	case CommandRun:
		panic("not implemented")
	default:
		panic("no command specified")
	}
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func runDebug() error {
	// copy config.toml to config dir/hat/config.toml
	configTOMLPath, err := getConfigTOMLPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(configTOMLPath), 0755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	cmd(fmt.Sprintf("cp config.toml %s", configTOMLPath))

	// build
	cmd("go build -o hat.build .")
	// run
	cmd("./hat.build")
	return nil
}
