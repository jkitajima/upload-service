package httphandler

import (
	"log"
	"net/http"
	"upload/pkg/file"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (s *fileServer) handleFileDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "fileID")
		uuid, err := uuid.Parse(id)
		if err != nil {
			log.Println("failed to parse uuid")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := file.Delete(r.Context(), s.db, uuid); err != nil {
			log.Println("failed to delete requested file")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
