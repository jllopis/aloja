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

// Exemple per posar un logger per defecte i permetre fer-ne servir d'altres alternatius.
//package authboss
//
//import (
//	"log"
//	"os"
//)
//
//// DefaultLogger is a basic logger.
//type DefaultLogger log.Logger
//
//// NewDefaultLogger creates a logger to stdout.
//func NewDefaultLogger() *DefaultLogger {
//	return ((*DefaultLogger)(log.New(os.Stdout, "", log.LstdFlags)))
//}
//
//// Write writes to the internal logger.
//func (d *DefaultLogger) Write(b []byte) (int, error) {
//	((*log.Logger)(d)).Printf("%s", b)
//	return len(b), nil
//}
