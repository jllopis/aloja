package aloja

import (
	"net/http"
	"path"
	"strings"

	"github.com/jllopis/aloja/mw"
	"github.com/julienschmidt/httprouter"
	"github.com/nbio/httpcontext"
)

type Subrouter struct {
	router *httprouter.Router
	*mw.Stack
	path string
}

func (s *Subrouter) NewSubrouter(subpath string, middlewares ...mw.Middleware) *Subrouter {
	sr := &Subrouter{router: s.router, path: subpath}
	sr.Stack = mw.New()
	return sr
}

type key int

const paramsKey key = 0

// Handle serves an endpoint with the provided handler
func (s *Subrouter) Handle(method string, path string, h http.Handler) {
	// calcular path
	fullPath := s.getFullPath(path)
	hrh := func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		httpcontext.Set(req, paramsKey, params)
		s.Stack.Then(h).ServeHTTP(w, req)
	}
	s.router.Handle(method, fullPath, hrh)
}

// HandleFunc serves an endpoint with the provided handler
func (s *Subrouter) HandleFunc(method string, path string, f func(w http.ResponseWriter, r *http.Request)) {
	s.Handle(method, path, http.HandlerFunc(f))
}

// Params returns the httprouter.Params for req.
func Params(req *http.Request) httprouter.Params {
	if value, ok := httpcontext.GetOk(req, paramsKey); ok {
		if params, ok := value.(httprouter.Params); ok {
			return params
		}
	}
	return httprouter.Params{}
}

func (s *Subrouter) getFullPath(p string) string {
	if p == "" {
		return s.path
	}
	fullPath := path.Join(s.path, p)
	if appendSlash := strings.HasSuffix(p, "/") && !strings.HasSuffix(fullPath, "/"); appendSlash {
		return fullPath + "/"
	}
	return fullPath
}

// Get registers a GET handler for the given path.
func (s *Subrouter) Get(path string, handler http.Handler) { s.Handle("GET", path, handler) }

// Put registers a PUT handler for the given path.
func (s *Subrouter) Put(path string, handler http.Handler) { s.Handle("PUT", path, handler) }

// Post registers a POST handler for the given path.
func (s *Subrouter) Post(path string, handler http.Handler) { s.Handle("POST", path, handler) }

// Patch registers a PATCH handler for the given path.
func (s *Subrouter) Patch(path string, handler http.Handler) { s.Handle("PATCH", path, handler) }

// Delete registers a DELETE handler for the given path.
func (s *Subrouter) Delete(path string, handler http.Handler) { s.Handle("DELETE", path, handler) }

// Options registers a OPTIONS handler for the given path.
func (s *Subrouter) Options(path string, handler http.Handler) { s.Handle("OPTIONS", path, handler) }

// ServeStatic provides a quick way to serve static files
func (s *Subrouter) ServeStatic(rpath string, dir string) {
	fp := s.getFullPath(rpath)
	fs := http.Dir(dir)
	fh := http.StripPrefix(fp, http.FileServer(fs))
	fp = path.Join(fp, "/*filepath")
	s.router.Handler("GET", fp, s.Stack.Then(fh))
	s.router.Handler("HEAD", fp, s.Stack.Then(fh))
}

func (s *Subrouter) Use(m ...mw.Middleware) {
	s.Stack.Add(m...)
}
