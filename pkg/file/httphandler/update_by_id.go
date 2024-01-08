package httphandler

import (
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
			http.Error(w, "failed to decode client input", http.StatusInternalServerError)
			return
		}

		id := chi.URLParam(r, "fileID")
		uuid, err := uuid.Parse(id)
		f.ID = uuid
		if err != nil {
			http.Error(w, "failed to parse uuid", http.StatusInternalServerError)
			return
		}

		if err := file.Update(r.Context(), s.db, uuid, &f); err != nil {
			http.Error(w, "failed to update requested file", http.StatusInternalServerError)
			return
		}

		if err := encoding.Respond(w, r, f, http.StatusOK); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
