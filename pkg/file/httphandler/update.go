package httphandler

import (
	"log"
	"net/http"
	"upload/pkg/file"
	"upload/util/encoding"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (s *fileServer) handleFileUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var f file.File

		if err := encoding.Decode(w, r, &f); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		id := chi.URLParam(r, "fileID")
		uuid, err := uuid.Parse(id)
		f.ID = uuid
		if err != nil {
			log.Println("failed to parse uuid")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := file.Update(r.Context(), s.db, uuid, &f); err != nil {
			log.Println("failed to update requested file")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := encoding.Respond(w, r, f, http.StatusCreated); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
