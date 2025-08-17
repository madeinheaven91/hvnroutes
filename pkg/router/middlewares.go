package router

import (
	"net/http"
	"slices"
)

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

type Middleware interface {
	// FIXME: how to change Handle to be of type MiddlewareFunc??
	Handle(http.HandlerFunc) http.HandlerFunc
}

type middlewareFuncWrapper struct {
	fn MiddlewareFunc
}

func (f middlewareFuncWrapper) Handle(next http.HandlerFunc) http.HandlerFunc {
	return f.fn(next)
}

func WrapMW(mw MiddlewareFunc) Middleware {
	return middlewareFuncWrapper{mw}
}

func (rt *route) buildChain() http.HandlerFunc {
	chain := func(handler http.HandlerFunc) http.HandlerFunc {
		for _, mw := range slices.Backward(rt.middlewareChain) {
			handler = mw.Handle(handler)
		}
		return handler
	}
	return chain(rt.handler.ServeHTTP)
}
