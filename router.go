package main

import (
	"net/http"
)

type IBackendRouter interface {
	UnprotectedHandle(path string, handler http.HandlerFunc)
	ProtectedHandle(path string, handler http.HandlerFunc)
	Use(mw IMiddleware)
	UseProtected(mw IMiddleware)
	http.Handler
}

type BackendRouter struct {
	mux *http.ServeMux

	globalMW    []IMiddleware
	protectedMW []IMiddleware
}

func NewBackendRouter() *BackendRouter {
	return &BackendRouter{
		mux: http.NewServeMux(),
	}
}

func (s *BackendRouter) UnprotectedHandle(path string, handler http.HandlerFunc) {
	finalHandler := s.applyMiddleware(handler, s.globalMW)
	s.mux.Handle(path, finalHandler)
}

func (s *BackendRouter) ProtectedHandle(path string, handler http.HandlerFunc) {
	finalHandler := s.applyMiddleware(handler, append(s.globalMW, s.protectedMW...))
	s.mux.Handle(path, finalHandler)
}

func (s *BackendRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *BackendRouter) Use(mw IMiddleware) {
	s.globalMW = append(s.globalMW, mw)
}

func (s *BackendRouter) UseProtected(mw IMiddleware) {
	s.protectedMW = append(s.protectedMW, mw)
}

func (s *BackendRouter) applyMiddleware(h http.Handler, mws []IMiddleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i].Wrap(h)
	}
	return h
}
