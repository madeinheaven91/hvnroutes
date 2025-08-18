package router_test

import (
	"io"
	"net/http"
	"testing"

	r "github.com/madeinheaven91/hvnroutes/pkg/router"
)

type test struct {
	name           string
	url            string
	fn             func(string) (resp *http.Response, err error)
	expectedStatus int
	expectedBody   string
}

func TestOmni(t *testing.T) {
	router := r.NewRouter(http.NewServeMux())
	router.Route = r.NewRoute("/root", nil, "GET").
		Route("/unspecified").
		Service(
			r.NewRoute("/test1").
				HandlerFunc(handlerFunc).
				Service(
					r.NewRoute("/{wildcard}").
						Handler(wildcardHandler{}),
				).
				Service(
					r.NewRoute("/mw").
						HandlerFunc(handlerFunc).
						MiddlewareFunc(middlewareFunc).
						MiddlewareFunc(middlewareFunc2),
				),
		)
	go http.ListenAndServe(":9000", router.BuildMux())

	tests := []test{
		{
			"root",
			"/root",
			http.Get,
			404,
			"",
		},
		// FIXME: when testing with curl, it blocks unallowed methods, but not here
		// {
		// 	"unallowed method",
		// 	"/root",
		// 	http.Head,
		// 	405,
		// 	"none",
		// },
		{
			"unspecified",
			"/root/unspecified",
			http.Get,
			404,
			"none",
		},
		{
			"not found",
			"/unknown",
			http.Get,
			404,
			"none",
		},
		{
			"test1",
			"/root/test1",
			http.Get,
			200,
			"test_handler",
		},
		{
			"middlewares",
			"/root/test1/mw",
			http.Get,
			200,
			"test_mw_1::test_handler::test_mw_2",
		},
		{
			"wildcard",
			"/root/test1/wildcard",
			http.Get,
			200,
			"wildcard",
		},
		{
			"wildcard2",
			"/root/test1/wildcard2",
			http.Get,
			200,
			"wildcard2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.fn("http://localhost:9000" + tt.url)
			if err != nil {
				panic(err)
			}
			if tt.expectedStatus != resp.StatusCode {
				t.Fatalf("status codes not maching. expected %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}
			if tt.expectedBody != "none" && string(body) != tt.expectedBody {
				t.Fatalf("bodies not maching. expected %s, got %s", tt.expectedBody, string(body))
			}
			resp.Body.Close()
		})
	}
}

type wildcardHandler struct{}

func (d wildcardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.PathValue("wildcard")))
}

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test_handler"))
}

func middlewareFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test_mw_1::"))
		next(w, r)
	}
}

func middlewareFunc2(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
		w.Write([]byte("::test_mw_2"))
	}
}
