package mw

import (
	"net/http"

	"github.com/rs/cors"
)

type CorsHandler cors.Cors

func NewCors(opts cors.Options) *CorsHandler {
	c := cors.New(opts)
	return (*CorsHandler)(c)
}

func (c *CorsHandler) Handler(h http.Handler) http.Handler {
	return (*cors.Cors)(c).Handler(h)
}
