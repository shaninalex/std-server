package main

import (
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/ory/x/prometheusx"
)

type RouterPublic struct {
	mux *http.ServeMux
	pmm *prometheusx.MetricsManager
}

type routerDeps interface {
	PrometheusManager() *prometheusx.MetricsManager
}

func NewRouterPublic(deps routerDeps) *RouterPublic {
	return &RouterPublic{
		mux: http.NewServeMux(),
		pmm: deps.PrometheusManager(),
	}
}

// NewTestRouterPublic creates a new RouterPublic for testing purposes without metrics.
func NewTestRouterPublic(*testing.T) *RouterPublic {
	return &RouterPublic{
		mux: http.NewServeMux(),
		pmm: nil, // No metrics manager in test environment
	}
}

func (r *RouterPublic) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r *RouterPublic) GET(path string, handler http.HandlerFunc) {
	r.HandlerFunc("GET", path, handler)
}

func (r *RouterPublic) HEAD(path string, handler http.HandlerFunc) {
	r.HandlerFunc("HEAD", path, handler)
}

func (r *RouterPublic) POST(path string, handler http.HandlerFunc) {
	r.HandlerFunc("POST", path, handler)
}

func (r *RouterPublic) PUT(path string, handler http.HandlerFunc) {
	r.HandlerFunc("PUT", path, handler)
}

func (r *RouterPublic) PATCH(path string, handler http.HandlerFunc) {
	r.HandlerFunc("PATCH", path, handler)
}

func (r *RouterPublic) DELETE(path string, handler http.HandlerFunc) {
	r.HandlerFunc("DELETE", path, handler)
}

func (r *RouterPublic) Handle(method, route string, handle http.HandlerFunc) {
	for _, pattern := range []string{
		method + " " + path.Join(route),
		method + " " + path.Join(route, "{$}"),
	} {
		handleWithAllMiddlewares(r.mux, r.pmm, pattern, handle)
	}
}

func (r *RouterPublic) HandlerFunc(method, route string, handler http.HandlerFunc) {
	for _, pattern := range []string{
		method + " " + path.Join(route),
		method + " " + path.Join(route, "{$}"),
	} {
		handleWithAllMiddlewares(r.mux, r.pmm, pattern, handler)
	}
}

func (r *RouterPublic) HandleFunc(pattern string, handler http.HandlerFunc) {
	for _, pattern := range []string{
		path.Join(pattern),
		path.Join(pattern, "{$}"),
	} {
		handleWithAllMiddlewares(r.mux, r.pmm, pattern, handler)
	}
}

func (r *RouterPublic) Handler(method, path string, handler http.Handler) {
	route := method + " " + path
	handleWithAllMiddlewares(r.mux, r.pmm, route, handler)
}

func (r *RouterPublic) HasRoute(method, path string) bool {
	_, pattern := r.mux.Handler(httptest.NewRequest(method, path, nil))
	return pattern != ""
}
