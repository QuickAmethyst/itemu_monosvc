package http

import (
	"github.com/rs/cors"
)

type CorsOptions = cors.Options

type Options struct {
	Cors *CorsOptions
}
