package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type HandleGetShortenedURLResponse struct {
	FullURL string `json:"full_url"`
}

func HandleGetShortenedURL(db map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")
		url, ok := db[code]
		if !ok {
			sendJSON(w, Response{Error: "not found"}, http.StatusNotFound)
			return
		}
		http.Redirect(w, r, url, http.StatusPermanentRedirect)
	}
}
