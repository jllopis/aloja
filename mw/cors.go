package mw

import (
	"net/http"

	"github.com/rs/cors"
)

type CorsOptions cors.Options

func CorsHandler(opts CorsOptions) Middleware {
	return func(h http.Handler) http.Handler {
		c := cors.New((cors.Options)(opts))
		return c.Handler(h)
	}
}
