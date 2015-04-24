package aloja

import (
	"net/http"
	"path"
	"strings"

	"github.com/dimfeld/httptreemux"
	"github.com/jllopis/aloja/mw"
	"github.com/nbio/httpcontext"
)

type Subrouter struct {
	router *httptreemux.TreeMux
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

type ParamCol map[string]string

func (p ParamCol) ByName(name string) string {
	if value, ok := p[name]; ok {
		return value
	}
	return ""
}

// Handle serves an endpoint with the provided handler
func (s *Subrouter) Handle(method string, path string, h http.Handler) {
	// calcular path
	fullPath := s.getFullPath(path)
	hrh := func(w http.ResponseWriter, req *http.Request, params map[string]string) {
		httpcontext.Set(req, paramsKey, params)
		s.Stack.Then(h).ServeHTTP(w, req)
	}
	s.router.Handle(method, fullPath, hrh)
}

// HandleFunc serves an endpoint with the provided handler
func (s *Subrouter) HandleFunc(method string, path string, f func(w http.ResponseWriter, r *http.Request)) {
	s.Handle(method, path, http.HandlerFunc(f))
}

// Params returns the httptreemux.Params for req.
func Params(req *http.Request) ParamCol {
	if value, ok := httpcontext.GetOk(req, paramsKey); ok {
		return ParamCol(value.(map[string]string))
	}
	return ParamCol{}
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
	root := path.Join(rpath, "/*filepath")
	// Add two more routes for handling index.html
	s.Handle("GET", rpath, fh)
	s.Handle("HEAD", rpath, fh)
	// and the files one...
	s.Handle("GET", root, fh)
	s.Handle("HEAD", root, fh)
}

func (s *Subrouter) Use(m ...mw.Middleware) {
	s.Stack.Add(m...)
}

// UseHandler registers an http.Handler as a middleware.
func (s *Subrouter) UseHandler(handler http.Handler) {
	s.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			handler.ServeHTTP(w, req)
			next.ServeHTTP(w, req)
		})
	})
}
