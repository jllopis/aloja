package mw

import (
	"bufio"
	"log"
	"net"
	"net/http"
)

// responseWriter is wrapper of http.ResponseWriter that records the HTTP status
// and the bytes written down
type responseWriter struct {
	w      http.ResponseWriter
	status int
	size   int
}

// Implements ResponseWriter interface
func (rw *responseWriter) Header() http.Header {
	return rw.w.Header()
}
func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.status == 0 {
		// The the first call to Write will trigger an implicit WriteHeader(http.StatusOK)
		// so we set status appropiately
		rw.status = http.StatusOK
	}
	size, err := rw.w.Write(b)
	rw.size += size
	return size, err
}
func (rw *responseWriter) WriteHeader(s int) {
	if rw.status > 0 {
		log.Printf("Headers already written!")
	}
	rw.w.WriteHeader(s)
	rw.status = s
}

// Implements http.Hijacker interface
func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	var conn net.Conn
	var r *bufio.ReadWriter
	var err error
	if conn, r, err = rw.w.(http.Hijacker).Hijack(); err == nil && rw.status == 0 {
		// If there is no error and Header has not been written, the status will be http.StatusSwitchingProtocols
		rw.status = http.StatusSwitchingProtocols
	}
	return conn, r, err
}

// Implements http.CloseNotifier interface
// Call the parent CloseNotify
func (rw *responseWriter) CloseNotify() <-chan bool {
	return rw.w.(http.CloseNotifier).CloseNotify()
}

// Implements the http.Flush interface
// Call the parent Flush
func (rw *responseWriter) Flush() {
	rw.w.(http.Flusher).Flush()
}

func (rw *responseWriter) Status() int {
	return rw.status
}
func (rw *responseWriter) Size() int {
	return rw.size
}
