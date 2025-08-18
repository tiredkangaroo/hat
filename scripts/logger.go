package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

// ColorHandler wraps an existing slog.Handler and adds color to log levels.
type ColorHandler struct {
	slog.Handler
}

func (h *ColorHandler) Handle(ctx context.Context, r slog.Record) error {
	t := r.Time.Format(time.RFC3339)
	switch r.Level {
	case slog.LevelError:
		r.Message = color.RedString(r.Message)
	case slog.LevelWarn:
		r.Message = color.YellowString(r.Message)
	}
	attrs := []string{}
	for a := range r.Attrs {
		attrs = append(attrs, fmt.Sprintf("%s=%v", a.Key, a.Value))
	}
	fmt.Fprintf(os.Stdout, "%s [%s] %s %s\n", t, r.Level, r.Message, strings.Join(attrs, " "))
	return nil
}

func initLogger() {
	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	// Wrap with ColorHandler
	logger := slog.New(&ColorHandler{Handler: textHandler})
	slog.SetDefault(logger)
}
