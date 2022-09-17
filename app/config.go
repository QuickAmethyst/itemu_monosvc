package main

import (
	"github.com/QuickAmethyst/monosvc/stdlibgo/grace"
	"github.com/QuickAmethyst/monosvc/stdlibgo/http"
	"github.com/QuickAmethyst/monosvc/stdlibgo/httpserver"
	"github.com/QuickAmethyst/monosvc/stdlibgo/sql"
	"github.com/go-redis/redis/v9"
)

type Config struct {
	Development        bool
	Grace              grace.Options
	Redis              redis.UniversalOptions
	AccountDatabase    sql.PostgresSQLOptions
	InventoryDatabase  sql.PostgresSQLOptions
	AccountingDatabase sql.PostgresSQLOptions
	HttpServer         httpserver.Options
	HttpCors           http.CorsOptions
}
