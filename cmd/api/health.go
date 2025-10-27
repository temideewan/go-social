package main

import (
	// "log"
	"net/http"
)

// CreatePost godoc
//
//	@Summary		Healthcheck
//	@Description	Ping to see if the server is still up
//	@Tags			ops
//	@Produce		json
//	@Success		200	{object}	string	"OK"
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}
	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}
}
