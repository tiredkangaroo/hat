package main

import (
	"log/slog"

	"github.com/tiredkangaroo/hat/proxy"
	"github.com/tiredkangaroo/hat/proxy/config"
)

func main() {
	if err := config.Init(); err != nil {
		slog.Error("initialize config", "error", err)
		return
	}
	slog.Info("configuration initialized")

	if err := proxy.Start(); err != nil {
		slog.Error("run proxy", "error", err)
		return
	}
}
