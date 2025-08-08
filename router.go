package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type IRouteHandlers interface {
	Get(path string, handler http.HandlerFunc)
	Post(path string, handler http.HandlerFunc)
	Use(middleware IMiddleware)
}

type IBackendRouter interface {
	Private() IRouteHandlers
	Public() IRouteHandlers
	Use(middleware IMiddleware)
	http.Handler
}

func NewBackendRouter() IBackendRouter {
	router := mux.NewRouter()

	return &BackendRouter{
		rootRouter: router,
		private:    &RouteHandlers{router: router},
		public:     &RouteHandlers{router: router},
	}
}

type BackendRouter struct {
	rootRouter *mux.Router
	private    *RouteHandlers
	public     *RouteHandlers
}

func (r *BackendRouter) Private() IRouteHandlers {
	return r.private
}

func (r *BackendRouter) Public() IRouteHandlers {
	return r.public
}

func (r *BackendRouter) Use(middleware IMiddleware) {
	r.public.middlewares = append(r.public.middlewares, middleware)
	r.private.middlewares = append(r.private.middlewares, middleware)
}

func (r *BackendRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.rootRouter.ServeHTTP(w, req)
}

type RouteHandlers struct {
	router      *mux.Router
	middlewares []IMiddleware
}

func (s *RouteHandlers) handle(path string, handler http.HandlerFunc, methods ...string) {
	h := http.Handler(handler)
	for i := len(s.middlewares) - 1; i >= 0; i-- {
		h = s.middlewares[i].Wrap(h)
	}
	s.router.Handle(path, h).Methods(methods...)
}

func (s *RouteHandlers) Get(path string, handler http.HandlerFunc) {
	s.handle(path, handler, http.MethodGet)
}

func (s *RouteHandlers) Post(path string, handler http.HandlerFunc) {
	s.handle(path, handler, http.MethodPost)
}

func (s *RouteHandlers) Use(middleware IMiddleware) {
	s.middlewares = append(s.middlewares, middleware)
}
