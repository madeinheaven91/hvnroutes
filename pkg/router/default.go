package router

import (
	"net/http"

	"github.com/madeinheaven91/hvnroutes/internal"
)

func Wrap(handlerFunc http.HandlerFunc) http.Handler {
	return &internal.HandlerFuncWrapper{
		HandlerFunc: handlerFunc,
	}
}
