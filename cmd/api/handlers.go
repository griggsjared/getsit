package main

import (
	"fmt"
	"net/http"
)

// rootHandler is the handler for the root path of the api.
func (a *app) rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", r.URL.Path)
}
