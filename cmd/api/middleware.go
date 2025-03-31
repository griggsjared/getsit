package main

import (
	"log/slog"
	"net/http"
	"time"
)

// middleware is a type that wraps an http.Handler and returns a new http.Handler
type middleware func(http.Handler) http.Handler

// middlewareStack will take a handler and a list of middlewares and return a new handler
func (a *app) middlewareStack(h http.Handler, middlewares ...middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

// middlewareStackFunc will take a handler and a list of middlewares and return a new handler.
// This is a convenience function that wraps middlewareStack to work with http.HandlerFunc
func (a *app) middlewareStackFunc(h http.HandlerFunc, middleware ...middleware) http.HandlerFunc {
	return a.middlewareStack(h, middleware...).ServeHTTP
}

// loggerMiddleware records the and method of the incoming request and the time taken to process the request
func (a *app) loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		ip := r.RemoteAddr
		if r.Header.Get("X-Forwarded-For") != "" {
			ip = r.Header.Get("X-Forwarded-For")
		}
		a.logger.Info("request", slog.String("ip", ip), slog.String("path", r.URL.Path), slog.String("method", r.Method), slog.Duration("duration", time.Since(start)))

	})
}
