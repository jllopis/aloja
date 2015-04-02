package mw

import (
	"net/http"

	"github.com/justinas/alice"
)

// Stack defines the method to stack middlewares to our running stack
type Stack struct {
	Stack alice.Chain
}

// Middleware is our middelware definition (standandard)
type Middleware func(http.Handler) http.Handler

// New get a handler for a new Middleware stack
func New() *Stack {
	return &Stack{Stack: alice.New()}
}

// Add adds a new middleware to the stack
func (s *Stack) Add(m ...Middleware) {
	for _, mdw := range m {
		s.Stack = s.Stack.Append(alice.Constructor(mdw))
	}
}

func (s *Stack) Then(h http.Handler) http.Handler {
	return s.Stack.Then(h)
}
