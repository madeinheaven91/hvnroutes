// Package router
//
// Package that provides router and middleware structs
package router

import (
	"net/http"
)

type Router struct {
	Mux   *http.ServeMux
	Route *route
}

// NewRouter creates new router
func NewRouter(mux *http.ServeMux) *Router {
	return &Router{
		Mux: mux,
	}
}

// BuildMux assembles all routes, middlewares and handlers and returns *http.ServeMux
func (c *Router) BuildMux() *http.ServeMux {
	c.Route.buildSubroutes(c.Mux, "")
	return c.Mux
}

type route struct {
	path            string
	handler         http.Handler
	method          string
	subroutes       []*route
	middlewareChain []Middleware
	stopProp        bool
}

// NewRoute сreates a new route. Leave `method` empty to allow all methods.
// If `handler` is nil, then it will be set to defaultHandler.
//
// params:
//   - handler http.Handler or http.HandlerFunc, default nil
//   - method: string, default ""
//   - middlewares []Middleware, default nil
//   - stopProp bool, default false
func NewRoute(path string, params ...any) *route {
	p := ParseParams(params...)
	return &route{
		path:            path,
		handler:         p.Handler,
		method:          p.Method,
		subroutes:       nil,
		middlewareChain: p.Middlewares,
		stopProp:        p.StopProp,
	}
}

// NewRouteWithParams сreates a new route with provided RouteParams struct. 
func NewRouteWithParams(path string, params RouteParams) *route {
	return &route{
		path:            path,
		handler:         params.Handler,
		method:          params.Method,
		subroutes:       nil,
		middlewareChain: params.Middlewares,
		stopProp:        params.StopProp,
	}
}

func (rt *route) buildSubroutes(mux *http.ServeMux, prefix string) {
	newPath := prefix + rt.path
	mux.Handle(rt.method+" "+newPath, rt.buildChain())

	for _, route := range rt.subroutes {
		// Stop middleware propagation if stopMiddleware is true or if current route has no middlewares
		if !route.stopProp || rt.middlewareChain == nil {
			route.middlewareChain = append(rt.middlewareChain, route.middlewareChain...)
		}
		go route.buildSubroutes(mux, newPath)
	}
}

// Path sets the path for this route.
// When Router is being assembled, this route's path is prepended with all it's parents' paths
func (rt *route) Path(path string) *route {
	rt.path = path
	return rt
}

// Handler sets the handler for this route
func (rt *route) Handler(handler http.Handler) *route {
	rt.handler = handler
	return rt
}

// HandlerFunc sets the handler function for this route
func (rt *route) HandlerFunc(handler http.HandlerFunc) *route {
	rt.handler = Wrap(handler)
	return rt
}

// Method adds a method to the allowed method's list. Don't use or set to "" to allow all methods.
func (rt *route) Method(method string) *route {
	rt.method = method
	return rt
}

// Middleware appends this routes' middleware chain with provided struct
func (rt *route) Middleware(mw Middleware) *route {
	rt.middlewareChain = append(rt.middlewareChain, mw)
	return rt
}

// MiddlewareFunc appends this routes' middleware chain with provided function
func (rt *route) MiddlewareFunc(mw MiddlewareFunc) *route {
	rt.middlewareChain = append(rt.middlewareChain, middlewareFuncWrapper{mw})
	return rt
}

// StopMiddleware stops middleware propagation to this route from it's parent
func (rt *route) StopMiddleware() *route {
	rt.stopProp = true
	return rt
}

// Route adds a new subroute to this route's children.
func (rt *route) Route(path string, params ...any) *route {
	p := ParseParams(params...)
	subroute := &route{
		path:            path,
		handler:         p.Handler,
		method:          p.Method,
		subroutes:       nil,
		middlewareChain: p.Middlewares,
		stopProp:        p.StopProp,
	}
	rt.subroutes = append(rt.subroutes, subroute)
	return rt
}

// RouteWithParams adds a new subroute to this route's children.
func (rt *route) RouteWithParams(path string, params ...any) *route {
	p := ParseParams(params...)
	subroute := &route{
		path:            path,
		handler:         p.Handler,
		method:          p.Method,
		subroutes:       nil,
		middlewareChain: p.Middlewares,
		stopProp:        p.StopProp,
	}
	rt.subroutes = append(rt.subroutes, subroute)
	return rt
}

// Service adds an existing subroute to this route's children.
func (rt *route) Service(service *route) *route {
	rt.subroutes = append(rt.subroutes, service)
	return rt
}
