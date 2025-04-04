package delete

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"url-shortener/internal/storage"
)

type URLDeleter interface {
	DeleteURL(alias string) error
}

type response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"` // Добавляем поле для сообщений об ошибках
}

func New(urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Printf("INFO: %s - alias is empty", op)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response{
				Status:  "error",
				Message: "alias is required",
			})
			return
		}

		err := urlDeleter.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Printf("INFO: %s - url not found, alias: %s", op, alias)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(response{
				Status:  "error",
				Message: "url not found",
			})
			return
		}
		if err != nil {
			log.Printf("ERROR: %s - failed to delete url: %v", op, err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response{
				Status:  "error",
				Message: "internal error",
			})
			return
		}

		log.Printf("INFO: %s - url deleted successfully, alias: %s", op, alias)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response{
			Status: "success",
		})
	}
}
