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

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func LoggingMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lrw := NewLoggingResponseWriter(w)

		start := time.Now()
		next(w, r)
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
