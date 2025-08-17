package proxy

import (
	"fmt"
	"log/slog"

	"net"

	"github.com/tiredkangaroo/hat/proxy/certificates"
	"github.com/tiredkangaroo/hat/proxy/config"

	"github.com/valyala/fasthttp"
)

type environment struct {
	listener    net.Listener
	certService *certificates.Service
}

var env *environment = &environment{}

func initialize() error {
	listener, err := net.Listen("tcp4", config.DefaultConfig.Addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", config.DefaultConfig.Addr, err)
	}
	slog.Info("listening on", "address", config.DefaultConfig.Addr)

	certService, err := certificates.GetService()
	if err != nil {
		return fmt.Errorf("initialize certificate service: %w", err)
	}

	env.listener = listener
	env.certService = certService
	return nil
}

func Start() error {
	defer env.listener.Close()

	if err := initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	}

	if err := fasthttp.Serve(env.listener, func(ctx *fasthttp.RequestCtx) {
		var err error
		if ctx.Method()[0] == 'C' { // CONNECT method (secure tunnel)
			err = handleHTTPS(ctx)
		} else { // HTTP proxy
			err = handleHTTP(ctx)
		}
		if err != nil {
			slog.Error("handle request", "error", err)
			ctx.SetStatusCode(fasthttp.StatusBadGateway)
		}
	}); err != nil {
		return fmt.Errorf("fasthttp listen and serve: %w", err)
	}

	return nil
}
