package aloja

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/fvbock/endless"
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

const (
	VERSION = "v0.1.0"
)

var (
	templates *template.Template
)

// New creates a new Aloja with options. Accepted are:
// - Host
// - Port
// - SSLConf
// It exposes a global middleware that is called on every request,
// independently of the sobrouter configured if any
func New(options ...func(s *Aloja) *Aloja) *Aloja {
	r := httprouter.New()
	r.HandleMethodNotAllowed = false
	srv := &Aloja{
		router:           r,
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
		log.Printf("Aloja %s started on %s:%s", VERSION, s.host, s.port)
		log.Fatalf("(ERR) main: Cannot start https server: %s", endless.ListenAndServeTLS(s.host+":"+s.port, s.cert, s.key, s.globalMiddleware.Then(s.router)))
	} else {
		// Non TLS available!!!
		log.Printf("Aloja %s started on %s:%s", VERSION, s.host, s.port)
		log.Printf("SSL disabled!! Please, provide a certificate to be secured!!")
		endless.ListenAndServe(s.host+":"+s.port, s.globalMiddleware.Then(s.router))
	}
}

// Add stacks a new middleware. It will be the last called.
// It accepts a list of middlewares and ordere is preserved left to right.
func (s *Aloja) AddGlobal(m ...mw.Middleware) {
	s.globalMiddleware.Add(m...)
}

func (s *Aloja) LoadTemplates(tdir string, templateDelims []string) (*template.Template, error) {
	// So sorry. Can't remember where I found this code. I'll give credit when find author.
	// initialize the templates,
	// couldn't have used http://golang.org/pkg/html/template/#ParseGlob
	// since we have custom delimiters.
	basePath := tdir
	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// don't process folders themselves
		if info.IsDir() {
			return nil
		}
		templateName := path[len(basePath):]
		if templates == nil {
			templates = template.New(templateName)
			templates.Delims(templateDelims[0], templateDelims[1])
			_, err = templates.ParseFiles(path)
		} else {
			_, err = templates.New(templateName).ParseFiles(path)
		}
		log.Printf("Processed template %s\n", templateName)
		return err
	})
	if err != nil {
		return nil, err
	}
	return templates, nil
}
