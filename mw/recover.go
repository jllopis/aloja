package mw

import (
	"fmt"
	"net/http"
	"runtime"
)

// Recover returns a middleware which recovers from panics anywhere in the chain
// and handles the control to the centralized HTTPErrorHandler.
func Recover(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("\n>>>> Called defer func on Recover mw\n")
				w.WriteHeader(http.StatusInternalServerError)
				stack := make([]byte, 1<<16)
				stack = stack[:runtime.Stack(stack, true)]
				fmt.Printf("PANIC: %s\n%s", err, stack)
				//fmt.Fprintf(w, "PANIC: %s\n%s", err, stack)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
