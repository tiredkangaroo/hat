package proxy

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"

	"github.com/valyala/fasthttp"
)

func handleHTTP(ctx *fasthttp.RequestCtx) error {
	return perform(&ctx.Request, &ctx.Response)
}

func handleHTTPS(ctx *fasthttp.RequestCtx) error {
	host := string(ctx.Host()) // string conversion because i do not want to mess with fasthttp memory management
	ctx.Hijack(func(c net.Conn) {
		defer c.Close()

		// send 200 ok response
		okResp := []byte("HTTP/1.1 200 OK\r\n" + "Content-Length: 0\r\n\r\n")
		if _, err := c.Write(okResp); err != nil {
			slog.Error("write ok response", "error", err)
			return
		}

		if env.certService.Enabled { // use mitm if enabled
			handleMITM(host, c)
			return
		}

		hostConn, err := net.Dial("tcp", host) // connect to the target server
		if err != nil {
			slog.Error("dial target server", "host", host, "error", err)
			return
		}
		defer hostConn.Close()

		wg := &sync.WaitGroup{}
		wg.Add(2)
		go func() {
			defer wg.Done()
			io.Copy(c, hostConn) // copy data from server to client
		}()
		go func() {
			defer wg.Done()
			io.Copy(hostConn, c) // copy data from client to server
		}()
		wg.Wait()
	})
	return nil
}

func handleMITM(host string, c net.Conn) error {
	tlsConn, err := env.certService.TLSConn(c, host)
	if err != nil {
		return fmt.Errorf("convert to TLS connection: %w", err)
	}
	defer tlsConn.Close()

	// we're not doing anything with the data yet, but it's here if needed in future
	fasthttp.ServeConn(tlsConn)
}

func perform(req *fasthttp.Request, resp *fasthttp.Response) error {
	req.Header.Del("Proxy-Authorization")
	req.Header.Del("Proxy-Connection")
	return fasthttp.Do(req, resp)
}
