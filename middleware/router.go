package middleware

import (
	"net/http"
)

type Middleware interface {
	Method()  MiddlewareHandler
	// Status returns whether the middleware is enabled
	Status()  bool
	// Experimental returns whether the middleware is experimental
	Experimental() bool
}

type MiddlewareHandler func(http.Handler) http.Handler