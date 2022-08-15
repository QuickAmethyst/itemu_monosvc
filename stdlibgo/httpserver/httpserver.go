package httpserver

import (
	"context"
	"github.com/QuickAmethyst/monosvc/stdlibgo/errors"
	"github.com/QuickAmethyst/monosvc/stdlibgo/http"
	"io"
	"log"
	"net"
	nativeHttp "net/http"
)

type HttpServer interface {
	Addr() string
	Serve(ln net.Listener) error
	Shutdown(ctx context.Context) error
}

type httpServer struct {
	opt    Options
	server *nativeHttp.Server
}

func (h *httpServer) Shutdown(ctx context.Context) error {
	return h.server.Shutdown(ctx)
}

func (h *httpServer) Addr() string {
	return h.opt.Address
}

func (h *httpServer) Serve(l net.Listener) error {
	if l == nil {
		return errors.PropagateWithCode(http.ErrListenerNil, errors.ErrHttpServerListenerNil, http.ErrListenerNil.Error())
	}

	return h.server.Serve(l)
}

func New(opt Options, handler nativeHttp.Handler, logWriter io.Writer) HttpServer {
	return &httpServer{
		opt: opt,
		server: &nativeHttp.Server{
			Addr:              opt.Address,
			Handler:           handler,
			ReadHeaderTimeout: opt.ReadHeaderTimeout,
			ReadTimeout:       opt.ReadTimeout,
			WriteTimeout:      opt.WriteTimeout,
			IdleTimeout:       opt.IdleTimeout,
			ErrorLog:          log.New(logWriter, "", log.LstdFlags),
		},
	}
}
