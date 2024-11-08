package main

import (
	"fmt"
	"net/http"
)

// rootHandler is the handler for the root path of the api.
func (a *app) rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", r.URL.Path)
}

// healthzHandler is the handler for the healthz path of the api.
func (a *app) healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
