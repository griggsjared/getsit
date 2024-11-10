package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/sessions"
)

const flashSessionName string = "flash-session"

// setFlashMessage sets a flash message in the session
func (a *app) setFlashMessage(w http.ResponseWriter, r *http.Request, message string) {
	session, err := a.session.Get(r, flashSessionName)
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}

	session.AddFlash(message, "message")
	session.Save(r, w)
}

// getFlashMessage gets a flash message from the session
func (a *app) getFlashMessage(w http.ResponseWriter, r *http.Request) string {

	session, err := a.session.Get(r, flashSessionName)
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return ""
	}

	defer session.Save(r, w)

	flashes := session.Flashes("message")
	if len(flashes) == 0 {
		return ""
	}

	return flashes[0].(string)
}

// setFlashErrors sets flash errors in the session
func (a *app) setFlashErrors(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	session, err := a.session.Get(r, flashSessionName)
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}

	if errors == nil {
		session.Save(r, w)
		return
	}

	errorsJson, err := json.Marshal(errors)
	if err != nil {
		http.Error(w, "Failed to encode errors", http.StatusInternalServerError)
		return
	}

	session.AddFlash(string(errorsJson), "errors")
	session.Save(r, w)
}

// getFlashErrors gets flash errors from the session
func (a *app) getFlashErrors(w http.ResponseWriter, r *http.Request) map[string]string {

	session, err := a.session.Get(r, flashSessionName)
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return nil
	}

	defer sessions.Save(r, w)

	flashes := session.Flashes("errors")
	if len(flashes) == 0 {
		return nil
	}

	errors := make(map[string]string)
	err = json.Unmarshal([]byte(flashes[0].(string)), &errors)
	if err != nil {
		http.Error(w, "Failed to decode errors", http.StatusInternalServerError)
		return nil
	}

	return errors
}

func (a *app) setFlashInputs(w http.ResponseWriter, r *http.Request, inputs map[string]string) {
	session, err := a.session.Get(r, flashSessionName)
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}

	if inputs == nil {
		session.Save(r, w)
		return
	}

	inputsJson, err := json.Marshal(inputs)
	if err != nil {
		http.Error(w, "Failed to encode inputs", http.StatusInternalServerError)
		return
	}

	session.AddFlash(string(inputsJson), "inputs")
	session.Save(r, w)
}

func (a *app) getFlashInputs(w http.ResponseWriter, r *http.Request) map[string]string {
	session, err := a.session.Get(r, flashSessionName)
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return nil
	}

	defer session.Save(r, w)

	flashes := session.Flashes("inputs")
	if len(flashes) == 0 {
		return nil
	}

	inputs := make(map[string]string)
	err = json.Unmarshal([]byte(flashes[0].(string)), &inputs)
	if err != nil {
		http.Error(w, "Failed to decode inputs", http.StatusInternalServerError)
		return nil
	}

	return inputs
}
