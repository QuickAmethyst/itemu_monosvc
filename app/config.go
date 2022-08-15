package main

import (
	"github.com/QuickAmethyst/monosvc/stdlibgo/grace"
	"github.com/QuickAmethyst/monosvc/stdlibgo/http"
	"github.com/QuickAmethyst/monosvc/stdlibgo/httpserver"
	"github.com/QuickAmethyst/monosvc/stdlibgo/sql"
)

type Config struct {
	Development       bool
	Grace             grace.Options
	InventoryDatabase sql.PostgresSQLOptions
	HttpServer        httpserver.Options
	HttpCors          http.CorsOptions
}
