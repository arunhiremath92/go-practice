package api

import (
	"fmt"
	"net/http"
)

type UrlRequestObj struct {
	UrlObject string `json:"url"`
}

func (app *App) ShortenUrlHandler(w http.ResponseWriter, r *http.Request) {

	var urlobject UrlRequestObj
	err := decodeJSON(w, r, &urlobject)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("failed to decode the request %s", err))
		return
	}
	shortenedUrl, err := app.config.UrlService.ShortenUrl(urlobject.UrlObject)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("failed to shorten the url %s", err))
		return
	}
	responseObj := map[string]string{
		"shortenedUrl": shortenedUrl,
	}
	writeJSON(w, http.StatusOK, responseObj)

}

func (app *App) RetrieveUrlHandler(w http.ResponseWriter, r *http.Request) {
	var urlobject UrlRequestObj
	err := decodeJSON(w, r, &urlobject)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to read incoming object")
		return
	}
	fullUrl, err := app.config.UrlService.GetFullUrl(urlobject.UrlObject)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("failed to retrieve the full url %s", err))
		return
	}

	responseObj := map[string]string{
		"shortenedUrl": urlobject.UrlObject,
		"fullUrl":      fullUrl,
	}
	writeJSON(w, http.StatusOK, responseObj)
}
