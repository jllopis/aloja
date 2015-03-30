package mw

import (
	"log"
	"net/http"
	"time"
)

// LogHandler logs the calls to the route
func LogHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		tIni := time.Now()
		rw := &responseWriter{w: w}
		next.ServeHTTP(rw, r)
		tEnd := time.Now()
		log.Printf("%v [%s] %q %v %v %v\n", r.RemoteAddr, r.Method, r.URL.String(), rw.Size(), rw.Status(), tEnd.Sub(tIni))
	}

	return http.HandlerFunc(fn)
}
