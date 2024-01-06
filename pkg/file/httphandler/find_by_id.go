package httphandler

import (
	"log"
	"net/http"

	"upload/pkg/file"
	"upload/util/encoding"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (s *fileServer) handleFileFindByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "fileID")

		uuid, err := uuid.Parse(id)
		if err != nil {
			log.Println("failed to parse uuid")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		f, err := file.FindByID(r.Context(), s.db, uuid)
		if err != nil {
			log.Println("could not find any file with provided id")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := encoding.Respond(w, r, f, http.StatusOK); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
