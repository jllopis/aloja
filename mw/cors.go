package mw

import (
	"net/http"

	"github.com/rs/cors"
)

func CorsHandler(opts cors.Options) Middleware {
	return func(h http.Handler) http.Handler {
		c := cors.New(opts)
		return c.Handler(h)
	}
}
