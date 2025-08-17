package router

import (
	"fmt"
	"net/http"

	"github.com/madeinheaven91/hvnroutes/internal"
)

// RouteParams defines params for a Route to pass in RouteWithParams function.
type RouteParams struct {
	Handler     http.Handler
	Method      string
	Middlewares []Middleware
	StopProp    bool
}

// DefaultParams returns the default params
func DefaultParams() *RouteParams {
	return &RouteParams{
		Handler:     internal.DefaultHandler{},
		Method:      "",
		Middlewares: nil,
		StopProp:    false,
	}
}

// ParseParams parses arguments in order of RouteParams stuct.
// If not all entries provided, then the default ones are used.
// If the length of params is >2, then the last argument is considered StopProp flag.
// All the entries between 1th and the last element are considered middlewares.
//
// params:
//   - handler http.Handler, default = defaultHandler
//   - method: string, default = ""
//   - middlewares []Middleware, default = nil
//   - stopProp bool, default = false
func ParseParams(params ...any) *RouteParams {
	res := DefaultParams()
	switch len(params) {
	case 0:
	default:
		stopProp, ok := params[len(params)-1].(bool)
		if !ok {
			panic("Invalid stopProp flag")
		}
		res.StopProp = stopProp

		middlewares := params[2 : len(params)-1]
		mws := make([]Middleware, len(middlewares))
		for i, middleware := range middlewares {
			mw, ok := middleware.(Middleware)
			if !ok {
				panic(fmt.Sprint("Invalid middleware in position ", i+2))
			}
			mws[i] = mw
		}
		if len(mws) != 0 {
			res.Middlewares = mws
		}
		fallthrough
	case 2:
		method, ok := params[1].(string)
		if !ok {
			panic("Invalid method")
		}
		switch method {
		case "":
		case http.MethodGet,
			http.MethodConnect, http.MethodDelete,
			http.MethodHead, http.MethodOptions,
			http.MethodPatch, http.MethodPost,
			http.MethodPut, http.MethodTrace:
			res.Method = method
		default:
			panic("Invalid method")
		}
		fallthrough
	case 1:
		if params[0] == nil {
			break
		}
		handler, ok := params[0].(http.Handler)
		if !ok {
			panic("Invalid handler")
		}
		res.Handler = handler
	}
	return res
}
