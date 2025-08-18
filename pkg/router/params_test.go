package router_test

import (
	"net/http"
	"testing"

	"github.com/madeinheaven91/hvnroutes/pkg/router"
)

var (
	tHandler http.Handler      = router.Wrap(handlerFunc)
	fstMw    router.Middleware = router.WrapMW(middlewareFunc)
	sndMw    router.Middleware = router.WrapMW(middlewareFunc2)
)

func TestParseParams(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		params     []any
		want       *router.RouteParams
		shouldFail bool
	}{
		{
			"empty params",
			[]any{},
			router.DefaultParams(),
			false,
		},
		{
			"only handler",
			[]any{tHandler},
			&router.RouteParams{tHandler, "", nil, false},
			false,
		},
		{
			"nil handler",
			[]any{nil},
			router.DefaultParams(),
			false,
		},
		{
			"all methods",
			[]any{tHandler, ""},
			&router.RouteParams{tHandler, "", nil, false},
			false,
		},
		{
			"method",
			[]any{tHandler, "GET"},
			&router.RouteParams{tHandler, "GET", nil, false},
			false,
		},
		{
			"invalid method",
			[]any{tHandler, "invalid"},
			&router.RouteParams{},
			true,
		},
		{
			"empty stopProp",
			[]any{tHandler, "GET", nil},
			&router.RouteParams{},
			true,
		},
		{
			"only stopProp",
			[]any{tHandler, "GET", true},
			&router.RouteParams{tHandler, "GET", nil, true},
			false,
		},
		{
			"nil middleware",
			[]any{tHandler, "GET", nil, true},
			&router.RouteParams{},
			true,
		},
		{
			"one middleware",
			[]any{tHandler, "GET", fstMw, true},
			&router.RouteParams{tHandler, "GET", []router.Middleware{fstMw}, true},
			false,
		},
		{
			"two middlewares",
			[]any{tHandler, "GET", fstMw, sndMw, true},
			&router.RouteParams{tHandler, "GET", []router.Middleware{fstMw, sndMw}, true},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if r != nil {
					if !tt.shouldFail {
						t.Errorf("Recovered. Error: %v", r)
					}
				} else {
					if tt.shouldFail {
						t.Errorf("Test unexpectedly didn't panic")
					}
				}
			}()
			got := router.ParseParams(tt.params...)
			if got.Handler != tt.want.Handler {
				t.Errorf("Handler = %v, want %v", got, tt.want)
			}
			if got.Method != tt.want.Method {
				t.Errorf("Method = %v, want %v", got, tt.want)
			}
			if len(got.Middlewares) != len(tt.want.Middlewares) {
				t.Errorf("Middlewares = %v, want %v", got, tt.want)
			}
			if got.StopProp != tt.want.StopProp {
				t.Errorf("StopProp = %v, want %v", got, tt.want)
			}
		})
	}
}
