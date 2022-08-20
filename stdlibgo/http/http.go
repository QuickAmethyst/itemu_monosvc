package http

import (
	"github.com/bmizerany/pat"
	"github.com/rs/cors"
	netHttp "net/http"
)

type Http interface {
	Handle(method Method, path string, fn netHttp.HandlerFunc)
	Handler() netHttp.Handler
}

type http struct {
	mux  *pat.PatternServeMux
	cors *cors.Cors
	opt  *Options
}

func (h *http) Handle(method Method, path string, handler netHttp.HandlerFunc) {
	handler = appendHeaderToContext(handler)

	h.mux.Add(string(method), path, handler)
}

func (h *http) Handler() netHttp.Handler {
	var handler netHttp.Handler
	if h.cors != nil {
		handler = h.cors.Handler(h.mux)
	}

	netHttp.Handle("/", handler)

	return handler
}

func New(opt Options) Http {
	result := http{
		mux: pat.New(),
		opt: &opt,
	}

	if opt.Cors != nil {
		result.cors = cors.New(*opt.Cors)
	}

	return &result
}
