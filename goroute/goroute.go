package goroute

import (
	"fmt"
	grmiddleware "go-route/middleware"
	"net/http"
	"strings"
)

type (
	Mux struct {
		BasePath       string
		RootMiddleware http.Handler
		ServeMux       *http.ServeMux
	}
	Route struct {
		RequestMethod string
		Path          string
		HandlerFunc   http.HandlerFunc
		Handler       http.Handler
		Middleware    []grmiddleware.Middleware
		SubRoutes     []Route
	}
)

func NewMux(basePath string) *Mux {
	return &Mux{
		BasePath: basePath,
		ServeMux: http.NewServeMux(),
	}
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if m.RootMiddleware != nil {
		m.RootMiddleware.ServeHTTP(w, r)
	} else {
		m.ServeMux.ServeHTTP(w, r)
	}
}

func (m *Mux) AddRootMiddleware(mw ...grmiddleware.Middleware) http.Handler {
	m.RootMiddleware = grmiddleware.ChainMiddleware(m, mw...)
	return m.RootMiddleware
}

func (m *Mux) AddRoute(route Route) {
	path := buildPath(route.RequestMethod, m.BasePath, route.Path)
	if len(route.SubRoutes) > 0 {
		m.AddSubRoutes(route.Path, route.SubRoutes, route.Middleware...)
		path = strings.TrimSuffix(path, "/")
	}

	var handler = route.Handler
	if route.HandlerFunc != nil {
		handler = route.HandlerFunc
	}

	if handler != nil {
		if len(route.Middleware) > 0 {
			handler = grmiddleware.ChainMiddleware(handler, route.Middleware...)
		}
		m.ServeMux.Handle(path, handler)
	}
}

func (m *Mux) AddRoutes(routes []Route, mw ...grmiddleware.Middleware) {
	for _, route := range routes {
		if len(mw) > 0 {
			route.Middleware = append(mw, route.Middleware...)
		}
		m.AddRoute(route)
	}
}

func (m *Mux) AddSubRoutes(path string, routes []Route, mw ...grmiddleware.Middleware) *Mux {
	subPath := buildPath("", m.BasePath, path)
	subMux := NewMux(subPath)
	subMux.AddRoutes(routes, mw...)
	m.ServeMux.Handle(subPath, subMux)
	return subMux
}

func buildPath(requestMethod string, base string, paths ...string) string {
	var path string
	for _, p := range paths {
		path = fmt.Sprintf("%s/%s", path, p)
	}
	path = fmt.Sprintf("%s /%s/%s", requestMethod, base, path)
	return clean(path)
}

func clean(path string) string {
	path = strings.TrimSpace(path)
	for strings.Contains(path, "//") {
		path = strings.ReplaceAll(path, "//", "/")
	}
	return path
}
