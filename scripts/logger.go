package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/fatih/color"
)

// ColorHandler wraps an existing slog.Handler and adds color to log levels.
type ColorHandler struct {
	slog.Handler
}

func (h *ColorHandler) Handle(ctx context.Context, r slog.Record) error {
	switch r.Level {
	case slog.LevelError:
		r.Message = color.RedString(r.Message)
	case slog.LevelWarn:
		r.Message = color.YellowString(r.Message)
	}
	return h.Handler.Handle(ctx, r)
}

func initLogger() {
	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	// Wrap with ColorHandler
	logger := slog.New(&ColorHandler{Handler: textHandler})
	slog.SetDefault(logger)
}
