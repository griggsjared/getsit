package main

import (
	"net/http"
	"time"
)

type middleware func(http.Handler) http.Handler

// stack will take a handler and a list of middlewares and return a new handler
func (a *app) middlewareStack(h http.Handler, middlewares ...middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

// stackFunc will take a handler function and a list of middlewares and return a new handler
func (a *app) middlewareStackFunc(h http.HandlerFunc, middleware ...middleware) http.HandlerFunc {
	return a.middlewareStack(h, middleware...).ServeHTTP
}

func (a *app) loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		a.logger.Println(time.Since(start), r.Method, r.URL.Path)
	})
}
