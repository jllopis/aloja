package aloja

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jllopis/aloja/mw"
)

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You are on the about page.")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}

func TestServer(t *testing.T) {
	//server := New().Host("inscripciones.acb.info").Port("8887").SSLConf("certs/cert.pem", "certs/key.pem")
	server := New().Port("8888")
	server.Add(mw.CorsHandler, mw.CompressHandler)
	server.Add(mw.LogHandler)
	server.HandleFunc("GET", "/about", aboutHandler)
	server.Get("/", http.HandlerFunc(indexHandler))
	server.ServeStatic("/info", "static")

	server.Run()
}
