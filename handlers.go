package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/griggsjared/getsit/internal/entity"
)

type UrlEntryRepository interface {
	// Save will url entry to the store
	Save(ctx context.Context, url entity.Url) (entry *entity.UrlEntry, new bool, err error)
	// SaveVisit will increment the number of times the url has been visited
	SaveVisit(ctx context.Context, token string) error
	// GetFromToken will get the url entry from the token
	GetFromToken(ctx context.Context, token string) (*entity.UrlEntry, error)
	// GetFromUrl will get the url entry from the url
	GetFromUrl(ctx context.Context, url string) (*entity.UrlEntry, error)
}

// appRouter is the router for the application
// it contains the necessary dependencies for the application routes
type appRouter struct {
	repo UrlEntryRepository
}

// setup will setup the routes for the application router
func (ar *appRouter) setup(mux *http.ServeMux) {
	mux.HandleFunc("GET /{$}", ar.handleHomepage)
	mux.HandleFunc("POST /create", ar.handleCreate)
	mux.HandleFunc("GET /i/{token}", ar.handleInfo)
	mux.HandleFunc("GET /{token}", ar.handleRedirect)
	mux.HandleFunc("/", ar.handleNotFound)
}

// handleHomepage will show the homepage of the application
// for now, this shows instructions on how to use the application
func (ar *appRouter) handleHomepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to Getsit")
	fmt.Fprintln(w, "To create a new short url, send a POST request to /create with the url parameter")
	fmt.Fprintln(w, "To get information about a short url, send a GET request to /i/{token}")
	fmt.Fprintln(w, "To redirect to a short url, send a GET request to /{token}")
}

// handleCreate will create a new short url from the long url
// The long url is sent as a POST request to /create
// if successful, we will redirect to /i/{token} to show the information about the url entry
func (ar *appRouter) handleCreate(w http.ResponseWriter, r *http.Request) {

	errors := make(map[string]error)

	url := entity.Url(r.FormValue("url"))
	if err := url.Validate(); err != nil {
		errors["url"] = err
	}

	if len(errors) > 0 {
		fmt.Fprintln(w, errors)
		return
	}

	entry, _, err := ar.repo.Save(r.Context(), url)
	if err != nil {
		fmt.Fprintln(w, "Error saving url")
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/i/%s", entry.Token), http.StatusMovedPermanently)
}

// handleRedirect will redirect to the long url from the short url
// The short url contains the token that is used to access the long url
// if successful, we record the visit and redirect to the long url
func (ar *appRouter) handleRedirect(w http.ResponseWriter, r *http.Request) {

	token := entity.UrlToken(r.PathValue("token"))
	if err := token.Validate(); err != nil {
		ar.handleNotFound(w, r)
		return
	}

	entry, err := ar.repo.GetFromToken(r.Context(), string(token))
	if err != nil {
		ar.handleNotFound(w, r)
		return
	}

	err = ar.repo.SaveVisit(r.Context(), string(entry.Token))
	if err != nil {
		fmt.Fprintln(w, "Error saving visit")
		return
	}

	http.Redirect(w, r, string(entry.Url), http.StatusFound)
}

// handleInfo will show the information about the url entry
// The token is sent as a GET request to /i/{token}
// if successful, we will show the url, token, and the number of times the url has been visited
func (ar *appRouter) handleInfo(w http.ResponseWriter, r *http.Request) {

	token := entity.UrlToken(r.PathValue("token"))
	if err := token.Validate(); err != nil {
		ar.handleNotFound(w, r)
		return
	}

	entry, err := ar.repo.GetFromToken(r.Context(), string(token))
	if err != nil {
		ar.handleNotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Url: %s\n", entry.Url)
	fmt.Fprintf(w, "Token: %s\n", entry.Token)
	fmt.Fprintf(w, "Visit Count: %d\n", entry.VisitCount)
}

// handleNotFound will show a 404 error message
// this is the default handler for when a route is not found and
// can be used to show return a 404 status from within other handlers
func (ar *appRouter) handleNotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404: Sorry not found", http.StatusNotFound)
}
