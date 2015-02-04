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

func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello little grasshoper!")
}

func TestServer(t *testing.T) {
	//server := New().Host("inscripciones.acb.info").Port("8887").SSLConf("certs/cert.pem", "certs/key.pem")
	// New server on port 8888
	server := New().Port("8888")

	// Use CORS Handler in every request
	server.AddGlobal(mw.CorsHandler)

	// Log every request
	server.AddGlobal(mw.LogHandler)

	// handler about on main route
	server.HandleFunc("GET", "/about", aboutHandler)

	// Create a subrouter on branch /v1
	r1 := server.NewSubrouter("/v1")

	// Accept compression on this subrouter
	r1.Use(mw.CompressHandler)

	// Give this route an index
	r1.Get("/", http.HandlerFunc(indexHandler))

	// Serve also some static content
	r1.ServeStatic("/info", "static")

	// Create another subrouter on /v2
	r2 := server.NewSubrouter("/v2")

	// Say hello!
	r2.Get("/hello", http.HandlerFunc(sayHello))

	server.Run()
}
