package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.errorLog.Printf("internal server error: %s\t path: %s\t error: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorLog.Printf("bad request error: %s\t path: %s\t error: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusBadRequest, err.Error())
}
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorLog.Printf("not found error: %s\t path: %s\t error: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusNotFound, "not found")
}
