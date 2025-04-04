package redirect

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"url-shortener/internal/storage"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

type errorResponse struct {
	Error string `json:"error"`
}

func New(urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Printf("INFO: %s - alias is empty", op)
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(errorResponse{Error: "invalid request"})
			if err != nil {
				return
			}
			return
		}

		resURL, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Printf("INFO: %s - url not found, alias: %s", op, alias)
			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode(errorResponse{Error: "not found"})
			if err != nil {
				return
			}
			return
		}
		if err != nil {
			log.Printf("ERROR: %s - failed to get url: %v", op, err)
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(errorResponse{Error: "internal error"})
			if err != nil {
				return
			}
			return
		}

		log.Printf("INFO: %s - got url: %s", op, resURL)
		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
