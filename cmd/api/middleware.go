package main

import (
	"net/http"
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
