package main

import (
	"encoding/json"
	"net/http"

	"github.com/griggsjared/getsit/internal/url"
)

// urlEntryResponse is the response struct for the url entry
type urlEntryResponse struct {
	Token      string `json:"token"`
	Url        string `json:"url"`
	VisitCount int    `json:"visit_count"`
}

// createUrlEntryHandler is the handler to create a new url entry
func (a *app) createUrlEntryHandler(w http.ResponseWriter, r *http.Request) {

	if exists, _ := a.urlService.GetUrlByUrl(r.Context(), &url.GetUrlByUrlInput{
		Url: r.FormValue("url"),
	}); exists != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(urlEntryResponse{
			Token:      exists.Token.String(),
			Url:        exists.Url.String(),
			VisitCount: exists.VisitCount,
		})
	}

	input := &url.SaveUrlInput{
		Url: r.FormValue("url"),
	}

	entry, err := a.urlService.SaveUrl(r.Context(), input)
	if err != nil {
		a.errorHandler(w, r, http.StatusBadRequest, "Failed to save url")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(urlEntryResponse{
		Token:      entry.Token.String(),
		Url:        entry.Url.String(),
		VisitCount: entry.VisitCount,
	})
}

// getUrlEntryHandler is the handler to get a single url entry by token
func (a *app) getUrlEntryHandler(w http.ResponseWriter, r *http.Request) {

	entry, err := a.urlService.GetUrlByToken(r.Context(), &url.GetUrlByTokenInput{
		Token: r.PathValue("token"),
	})
	if err != nil {
		a.errorHandler(w, r, http.StatusNotFound, "Url entry not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(urlEntryResponse{
		Token:      entry.Token.String(),
		Url:        entry.Url.String(),
		VisitCount: entry.VisitCount,
	})
}

// errorResponse is the response struct for errors
type errorResponse struct {
	Message string `json:"message"`
}

// errorHandler is the handler for errors
func (a *app) errorHandler(w http.ResponseWriter, _ *http.Request, status int, message string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(errorResponse{
		Message: message,
	})
}

// healthzHandler is the handler for the healthz path of the api.
func (a *app) healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
