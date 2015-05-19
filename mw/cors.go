package mw

import (
	"net/http"

	"github.com/rs/cors"
)

func CorsHandler(h http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"localhost", "acb.info"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "OPTIONS", "PATCH", "DELETE"},
		AllowCredentials: true,
	})
	return c.Handler(h)
	//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		// Allow cross domain Ajax requests
	//		//if r.Method == "OPTIONS" {
	//		if origin := r.Header.Get("Origin"); origin != "" {
	//			w.Header().Set("Access-Control-Allow-Origin", origin)
	//		}
	//		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
	//		w.Header().Set("Access-Control-Request-Method", "*")
	//		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Origin, X-Requested-With, Accept, Authorization")
	//		w.Header().Set("Access-Control-Allow-Credentials", "true")
	//		//}
	//		//w.Header().Set("Access-Control-Allow-Origin", "*")
	//		h.ServeHTTP(w, r)
	//	})
}
