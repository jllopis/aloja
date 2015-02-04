package aloja

import (
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/jllopis/aloja/mw"
	"github.com/julienschmidt/httprouter"
)

// Aloja has the methods to abstract the work away. It provides a default Subrouter to
// carry on requests to '/'.
// You can coall NewSubrouter method on it to group routes and apply different middlewares to them
type Aloja struct {
	*Subrouter
	router           *httprouter.Router
	globalMiddleware *mw.Stack
	cert, key        string
	host             string
	port             string
}

// New creates a new Aloja with options. Accepted are:
// - Host
// - Port
// - SSLConf
// It exposes a global middleware that is called on every request,
// independently of the sobrouter configured if any
func New(options ...func(s *Aloja) *Aloja) *Aloja {
	srv := &Aloja{
		router:           httprouter.New(),
		globalMiddleware: mw.New(),
		host:             "",
		port:             "8888",
	}
	if options != nil {
		for _, option := range options {
			option(srv)
		}
	}
	srv.Subrouter = &Subrouter{
		srv.router,
		mw.New(),
		"/",
	}

	// add default service to give server time
	srv.HandleFunc("GET", "/time", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		t, _ := time.Now().UTC().MarshalJSON()
		w.Write(t)
	})
	return srv
}

// Host set the hostname for the running server
func (s *Aloja) Host(h string) *Aloja {
	s.host = h
	return s
}

// Port set the port for the running server
func (s *Aloja) Port(p string) *Aloja {
	s.port = p
	return s
}

// SSLConf pass in the cert and key for StartTLS
func (s *Aloja) SSLConf(cert, key string) *Aloja {
	if cert != "" && key != "" {
		s.cert = cert
		s.key = key
	}
	return s
}

// Run is a convenience function to start an http server.
// If cert and key were provided, ONLY https will be started!
// If they aren't, http will boot up. NOT RECOMMENDED!
func (s *Aloja) Run() {
	if s.cert != "" && s.key != "" {
		// StartTLS
		glog.Infof("Aloja started on %s:%s", s.host, s.port)
		glog.Fatalf("(ERR) main: Cannot start https server: %s", http.ListenAndServeTLS(s.host+":"+s.port, s.cert, s.key, s.globalMiddleware.Then(s.router)))
	} else {
		// Non TLS available!!!
		glog.Infof("Aloja started on %s:%s", s.host, s.port)
		glog.Warningf("SSL disabled!! Please, provide a certificate to be secured!!")
		http.ListenAndServe(s.host+":"+s.port, s.globalMiddleware.Then(s.router))
	}
}

// Add stacks a new middleware. It will be the last called.
// It accepts a list of middlewares and ordere is preserved left to right.
func (s *Aloja) AddGlobal(m ...mw.Middleware) {
	s.globalMiddleware.Add(m...)
}
