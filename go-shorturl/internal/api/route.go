package api

import (
	"net/http"
)

func (app *App) registerRoutesv1(mux *http.ServeMux) {
	mux.HandleFunc("POST /shorten", app.ShortenUrlHandler)
	mux.HandleFunc("GET /full", app.RetrieveUrlHandler)
}
