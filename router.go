package main

import (
	"net/http"
	"path"
)

type Middleware func(http.Handler) http.Handler

type Router struct {
	mux         *http.ServeMux
	middlewares []Middleware
}

func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r *Router) GET(path string, handler http.HandlerFunc) {
	r.HandlerFunc("GET", path, handler)
}

func (r *Router) HEAD(path string, handler http.HandlerFunc) {
	r.HandlerFunc("HEAD", path, handler)
}

func (r *Router) POST(path string, handler http.HandlerFunc) {
	r.HandlerFunc("POST", path, handler)
}

func (r *Router) PUT(path string, handler http.HandlerFunc) {
	r.HandlerFunc("PUT", path, handler)
}

func (r *Router) PATCH(path string, handler http.HandlerFunc) {
	r.HandlerFunc("PATCH", path, handler)
}

func (r *Router) DELETE(path string, handler http.HandlerFunc) {
	r.HandlerFunc("DELETE", path, handler)
}

func (r *Router) HandlerFunc(method, route string, handler http.HandlerFunc) {
	for _, pattern := range []string{
		method + " " + path.Join(route),
		method + " " + path.Join(route, "{$}"),
	} {
		r.handleWithAllMiddlewares(r.mux, pattern, handler)
	}
}

func (r *Router) Use(m Middleware) {
	r.middlewares = append(r.middlewares, m)
}

func (r *Router) handleWithAllMiddlewares(mux *http.ServeMux, pattern string, handler http.Handler) {
	// Apply global middlewares
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}

	mux.HandleFunc(pattern, func(w http.ResponseWriter, req *http.Request) {
		//NoCache(w)
		handler.ServeHTTP(w, req)
	})
}
