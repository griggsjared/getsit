package main

import (
	"fmt"
	"net/http"

	"github.com/griggsjared/getsit/internal/url"
	"github.com/griggsjared/getsit/web/template"

	"github.com/gorilla/csrf"
)

// homepageHandler will show the homepage of the application that shows the form to create a new short url
func (a *app) homepageHandler(w http.ResponseWriter, r *http.Request) {

	token := csrf.Token(r)

	err := template.Homepage(template.HomepageViewModel{
		CsrfToken: token,
		Message:   a.getFlashMessage(w, r),
		Errors:    a.getFlashErrors(w, r),
		Inputs:    a.getFlashInputs(w, r),
	}).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render homepage", http.StatusInternalServerError)
		return
	}
}

// createHandler will create a new short url from the long url
// The long url is sent as a POST request to /create
// if successful, we will redirect to /i/{token} to show the information about the url entry
func (a *app) createHandler(w http.ResponseWriter, r *http.Request) {

	if exists, _ := a.urlService.GetUrlByUrl(r.Context(), &url.GetUrlByUrlInput{
		Url: r.FormValue("url"),
	}); exists != nil {
		http.Redirect(w, r, fmt.Sprintf("/i/%s", exists.Token), http.StatusMovedPermanently)
		return
	}

	input := &url.SaveUrlInput{
		Url: r.FormValue("url"),
	}

	entry, err := a.urlService.SaveUrl(r.Context(), input)
	if err != nil {
		if len(input.ValidationErrors) > 0 {
			a.setFlashErrors(w, r, input.ValidationErrors)
		} else {
			a.setFlashErrors(w, r, map[string]string{"error": "Failed to save url"})
		}
		a.setFlashInputs(w, r, map[string]string{"url": r.FormValue("url")})
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/i/%s", entry.Token), http.StatusMovedPermanently)
}

// redirectHandler will redirect to the long url from the short url
// The short url contains the token that is used to access the long url
// if successful, we record the visit and redirect to the long url
func (a *app) redirectHandler(w http.ResponseWriter, r *http.Request) {

	entry, err := a.urlService.GetUrlByToken(r.Context(), &url.GetUrlByTokenInput{
		Token: r.PathValue("token"),
	})
	if err != nil {
		a.notFoundHandler(w, r)
		return
	}

	err = a.urlService.VisitUrlByToken(r.Context(), &url.VisitUrlByTokenInput{
		Token: entry.Token.String(),
	})
	if err != nil {
		fmt.Fprintln(w, "Error saving visit")
		return
	}

	http.Redirect(w, r, entry.Url.String(), http.StatusFound)
}

// infoHandler will show the information about the url entry
// The token is sent as a GET request to /i/{token}
// if successful, we will show the url, token, and the number of times the url has been visited
func (a *app) infoHandler(w http.ResponseWriter, r *http.Request) {

	entry, err := a.urlService.GetUrlByToken(r.Context(), &url.GetUrlByTokenInput{
		Token: r.PathValue("token"),
	})
	if err != nil {
		a.notFoundHandler(w, r)
		return
	}

	proto := "http"
	if r.TLS != nil {
		proto = "https"
	}

	err = template.Info(template.InfoViewModel{
		ShortUrl:          fmt.Sprintf("%s/%s", r.Host, entry.Token),
		ShortUrlWithProto: fmt.Sprintf("%s://%s/%s", proto, r.Host, entry.Token),
		Url:               entry.Url.String(),
		Token:             entry.Token.String(),
		VisitCount:        entry.VisitCount,
	}).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render information page", http.StatusInternalServerError)
		return
	}
}

// notFoundHandler will show a 404 error message
// this is the default handler for when a route is not found and
// can be used to show return a 404 status from within other handlers
func (a *app) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	err := template.ServerError(template.ServerErrorViewModel{
		Code: http.StatusNotFound,
		Msg:  "404: Page not found",
		Desc: "Sorry we could not find what you were looking for.",
	}).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render 404 page", http.StatusInternalServerError)
		return
	}
}

// tokenMismatchHandler will show a 403 error message
// this is the handler for when a csrf token mismatch occurs
func (a *app) tokenMismatchHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	err := template.ServerError(template.ServerErrorViewModel{
		Code: http.StatusNotFound,
		Msg:  "Invalid Request",
		Desc: "Sorry, the request was invalid. Please go back and try again.",
	}).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render the token mismatch page", http.StatusInternalServerError)
		return
	}
}

// healthzHandler is the handler for the healthz path of the web application.
func (a *app) healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
