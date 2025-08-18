package router

import (
	"net/http"

	"github.com/madeinheaven91/hvnroutes/internal"
)

func WrapHandler(handlerFunc http.HandlerFunc) http.Handler {
	return &internal.HandlerFuncWrapper{
		HandlerFunc: handlerFunc,
	}
}

func WrapMW(mw MiddlewareFunc) Middleware {
	return middlewareFuncWrapper{mw}
}
