package aloja

import (
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/jllopis/aloja/mw"
	"github.com/julienschmidt/httprouter"
)

// Server has the methods to abstract the work away
type Server struct {
	router    *httprouter.Router
	mwStack   *mw.Stack
	cert, key string
	host      string
	port      string
}

// New creates a new Server with options
func New(options ...func(s *Server) *Server) *Server {
	srv := &Server{
		router:  httprouter.New(),
		mwStack: mw.New(),
		host:    "",
		port:    "8888",
	}
	if options != nil {
		for _, option := range options {
			option(srv)
		}
	}
	// add default service to give server time
	srv.router.HandlerFunc("GET", "/time", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		t, _ := time.Now().UTC().MarshalJSON()
		w.Write(t)
	})
	return srv
}

// Host set the hostname for the running server
func (s *Server) Host(h string) *Server {
	s.host = h
	return s
}

// Port set the port for the running server
func (s *Server) Port(p string) *Server {
	s.port = p
	return s
}

// SSLConf pass in the cert and key for StartTLS
func (s *Server) SSLConf(cert, key string) *Server {
	if cert != "" && key != "" {
		s.cert = cert
		s.key = key
	}
	return s
}

// Run is a convenience function to start an http server.
// If cert and key were provided, ONLY https will be started!
// If they aren't, http will boot up. NOT RECOMMENDED!
func (s *Server) Run() {
	if s.cert != "" && s.key != "" {
		// StartTLS
		glog.Infof("Server started on %s:%s", s.host, s.port)
		glog.Fatalf("(ERR) main: Cannot start https server: %s", http.ListenAndServeTLS(s.host+":"+s.port, s.cert, s.key, s.mwStack.Stack.Then(s.router)))
	} else {
		// Non TLS available!!!
		glog.Infof("Server started on %s:%s", s.host, s.port)
		glog.Warningf("SSL disabled!! Please, provide a certificate to be secured!!")
		http.ListenAndServe(s.host+":"+s.port, s.mwStack.Stack.Then(s.router))
	}
}

// Add stacks a new middleware. It will be the last called.
// It accepts a list of middlewares and ordere is preserved left to right.
func (s *Server) Add(m ...mw.Middleware) {
	s.mwStack.Add(m...)
}

// Handle serves an endpoint with the provided handler
func (s *Server) Handle(method string, path string, f http.Handler) {
	s.router.Handler(method, path, f)
}

// HandleFunc serves an endpoint with the provided handler
func (s *Server) HandleFunc(method string, path string, f func(w http.ResponseWriter, r *http.Request)) {
	s.router.HandlerFunc(method, path, f)
}

// Get registers a GET handler for the given path.
func (s *Server) Get(path string, handler http.Handler) { s.Handle("GET", path, handler) }

// Put registers a PUT handler for the given path.
func (s *Server) Put(path string, handler http.Handler) { s.Handle("PUT", path, handler) }

// Post registers a POST handler for the given path.
func (s *Server) Post(path string, handler http.Handler) { s.Handle("POST", path, handler) }

// Patch registers a PATCH handler for the given path.
func (s *Server) Patch(path string, handler http.Handler) { s.Handle("PATCH", path, handler) }

// Delete registers a DELETE handler for the given path.
func (s *Server) Delete(path string, handler http.Handler) { s.Handle("DELETE", path, handler) }

// Options registers a OPTIONS handler for the given path.
func (s *Server) Options(path string, handler http.Handler) { s.Handle("OPTIONS", path, handler) }

// ServeStatic provides a quick way to serve static files
func (s *Server) ServeStatic(path string, dir string) {
	fs := http.Dir(dir)
	fh := http.StripPrefix(path, http.FileServer(fs))
	s.router.Handler("GET", path, fh)
	s.router.Handler("HEAD", path, fh)
}
