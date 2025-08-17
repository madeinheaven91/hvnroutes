package internal

import "net/http"

type DefaultHandler struct{}

func (d DefaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
}

type HandlerFuncWrapper struct {
	HandlerFunc http.HandlerFunc
}

func (h HandlerFuncWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.HandlerFunc(w, r)
}
