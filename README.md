aloja - little opinionated framework for plain http.Handlers
============================================================

This framework provide convenience for the usual suspects when developing RESTful APIs or Web Applications.

It uses mainstream and common utilities and it is heavily opinionated with focus to alleviate repetitive tasks.

# Features

In this release, we have:

- Fast routing (by way of [httprouter](https://github.com/julienschmidt/httprouter)
- Per request Context (provided by [httpcontext](https://github.com/nbio/httpcontext))
- Middleware (gently managed by [alice](https://github.com/justinas/alice))
- Named params
- Serve static pages
- Serve http.Templates
- Some middleware:
-- compression
-- logs
-- cors

# Installation

```
go get github.com/jllopis/aloja
```

# Quick sample

The framework lacks documentation. The use is something like this:

```
package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jllopis/aloja"
	"github.com/jllopis/aloja/mw"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello little grasshoper!")
}

func sayHelloName(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello %s!", aloja.Params(r).ByName("name"))
}

func main() {
	// New server on port 8888
	server := New().Port("8888")

	// Use CORS Handler in every request
	server.AddGlobal(mw.CorsHandler)

	// Log every request
	server.AddGlobal(mw.LogHandler)

	// handler on main route
	server.HandleFunc("GET", "/", indexHandler)

	// Create a subrouter on branch /v1
	r1 := server.NewSubrouter("/v1")

	// Accept compression on this subrouter
	r1.Use(mw.CompressHandler)

	// Say hello!
	r1.Get("/hello", http.HandlerFunc(sayHello))

	// Say hello!
	r1.Get("/hello/:name", http.HandlerFunc(sayHelloName))

	// Run rabbit run!
	server.Run()
}
```
