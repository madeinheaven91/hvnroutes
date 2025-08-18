// Package middlewares
//
// Provides a set of predefined middlewares
package middlewares

import (
	"log"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func loggingMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lrw := newLoggingResponseWriter(w)

		start := time.Now()
		next(lrw, r)
		end := time.Now()

		log.Printf("%s %s %s | %d sent %d | handled in %d ms",
			r.RemoteAddr,
			r.Method,
			r.URL.String(),
			lrw.statusCode,
			r.ContentLength,
			end.Sub(start).Milliseconds())
	}
}

type Logging struct{}

func (l Logging) Handle(next http.HandlerFunc) http.HandlerFunc {
	return loggingMW(next)
}
