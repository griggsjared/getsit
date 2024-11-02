package main

import (
	"fmt"
	"net/http"

	"github.com/griggsjared/getsit/internal"
	"github.com/griggsjared/getsit/web/template"

	"github.com/gorilla/csrf"
)

// homepageHandler will show the homepage of the application
// for now, this shows instructions on how to use the application
func (a *app) homepageHandler(w http.ResponseWriter, r *http.Request) {

	token := csrf.Token(r)

	err := template.Homepage(template.HomepageViewModel{
		CsrfToken: token,
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

	input := &internal.SaveUrlInput{
		Url: r.FormValue("url"),
	}

	entry, err := a.service.SaveUrl(r.Context(), input)
	if err != nil {
		// TODO: we need to handle the errors here so they are not just dumped to the screen.
		// We can redirect back with an error message to show on the page
		if len(input.ValidationErrors) > 0 {
			fmt.Fprintln(w, "Validation Errors:")
			for _, v := range input.ValidationErrors {
				fmt.Fprintf(w, "%s", v)
			}
		} else {
			fmt.Fprintln(w, "Error saving url")
		}
		// http.Redirect(w, r, "/?err=Oops", http.StatusMovedPermanently)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/i/%s", entry.Token), http.StatusMovedPermanently)
}

// redirectHandler will redirect to the long url from the short url
// The short url contains the token that is used to access the long url
// if successful, we record the visit and redirect to the long url
func (a *app) redirectHandler(w http.ResponseWriter, r *http.Request) {

	entry, err := a.service.GetUrl(r.Context(), &internal.GetUrlInput{
		Token: r.PathValue("token"),
	})
	if err != nil {
		a.notFoundHandler(w, r)
		return
	}

	err = a.service.VisitUrl(r.Context(), &internal.VisitUrlInput{
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

	entry, err := a.service.GetUrl(r.Context(), &internal.GetUrlInput{
		Token: r.PathValue("token"),
	})
	if err != nil {
		a.notFoundHandler(w, r)
		return
	}

	err = template.Info(template.InfoViewModel{
		ShortUrl:   fmt.Sprintf("%s/%s", r.Host, entry.Token),
		Url:        entry.Url.String(),
		Token:      entry.Token.String(),
		VisitCount: entry.VisitCount,
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
	err := template.NotFound().Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render 404 page", http.StatusInternalServerError)
		return
	}
}
