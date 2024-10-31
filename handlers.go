package main

import (
	"fmt"
	"net/http"

	"github.com/griggsjared/getsit/internal"
	"github.com/griggsjared/getsit/web/template"
)

// appHandler is the router for the application
// it contains the necessary dependencies for the application routes
type appHandler struct {
	service *internal.Service
}

// setup will setup the routes for the application router
func (ah *appHandler) setup(mux *http.ServeMux) {
	mux.HandleFunc("GET /{$}", ah.handleHomepage)
	mux.HandleFunc("POST /create", ah.handleCreate)
	mux.HandleFunc("GET /i/{token}", ah.handleInfo)
	mux.HandleFunc("GET /{token}", ah.handleRedirect)
	mux.HandleFunc("/", ah.handleNotFound)
}

// handleHomepage will show the homepage of the application
// for now, this shows instructions on how to use the application
func (ah *appHandler) handleHomepage(w http.ResponseWriter, r *http.Request) {
	err := template.Homepage().Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render homepage", http.StatusInternalServerError)
		return
	}
}

// handleCreate will create a new short url from the long url
// The long url is sent as a POST request to /create
// if successful, we will redirect to /i/{token} to show the information about the url entry
func (ah *appHandler) handleCreate(w http.ResponseWriter, r *http.Request) {

	input := &internal.SaveUrlInput{
		Url: r.FormValue("url"),
	}

	entry, err := ah.service.SaveUrl(r.Context(), input)
	if err != nil {
		fmt.Fprintln(w, input.ValidationErrors)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/i/%s", entry.Token), http.StatusMovedPermanently)
}

// handleRedirect will redirect to the long url from the short url
// The short url contains the token that is used to access the long url
// if successful, we record the visit and redirect to the long url
func (ah *appHandler) handleRedirect(w http.ResponseWriter, r *http.Request) {

	entry, err := ah.service.GetUrl(r.Context(), &internal.GetUrlInput{
		Token: r.PathValue("token"),
	})
	if err != nil {
		ah.handleNotFound(w, r)
		return
	}

	err = ah.service.VisitUrl(r.Context(), &internal.VisitUrlInput{
		Token: entry.Token.String(),
	})
	if err != nil {
		fmt.Fprintln(w, "Error saving visit")
		return
	}

	http.Redirect(w, r, entry.Url.String(), http.StatusFound)
}

// handleInfo will show the information about the url entry
// The token is sent as a GET request to /i/{token}
// if successful, we will show the url, token, and the number of times the url has been visited
func (ah *appHandler) handleInfo(w http.ResponseWriter, r *http.Request) {

	entry, err := ah.service.GetUrl(r.Context(), &internal.GetUrlInput{
		Token: r.PathValue("token"),
	})
	if err != nil {
		ah.handleNotFound(w, r)
		return
	}

	err = template.Info(template.InfoViewModel{
		Url:        entry.Url.String(),
		Token:      entry.Token.String(),
		VisitCount: entry.VisitCount,
	}).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render information page", http.StatusInternalServerError)
		return
	}
}

// handleNotFound will show a 404 error message
// this is the default handler for when a route is not found and
// can be used to show return a 404 status from within other handlers
func (ah *appHandler) handleNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	err := template.NotFound().Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render 404 page", http.StatusInternalServerError)
		return
	}
}
