package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/griggsjared/getsit/web/template"
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

// templateColorMiddleware will set the color-mode context value based on the color-mode cookie
func (a *app) templateColorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mode := "dark"
		if colorModeCookie, err := r.Cookie("color-mode"); err == nil {
			if colorModeCookie.Value == "light" {
				mode = "light"
			}
		}
		ctx := context.WithValue(r.Context(), template.ColorModeCtxKey, mode)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
