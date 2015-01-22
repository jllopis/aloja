package mw

import (
	"log"
	"net/http"
	"time"
)

// LogHandler logs the calls to the route
// TODO: Afegir codi d'estat. Veure https://github.com/gocraft/web/blob/master/response_writer.go
// per implementar wrapper sobre http.ResponseWriter y capturar l'estat
func LogHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		tIni := time.Now()
		next.ServeHTTP(w, r)
		tEnd := time.Now()
		log.Printf("%v [%s] %q %v\n", r.RemoteAddr, r.Method, r.URL.String(), tEnd.Sub(tIni))
	}

	return http.HandlerFunc(fn)
}
