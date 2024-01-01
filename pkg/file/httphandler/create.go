package httphandler

import (
	"net/http"

	"upload/pkg/file"
	"upload/util/encoding"
)

func (s *fileServer) handleFileCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var f file.File

		if err := encoding.Decode(w, r, &f); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := file.Create(r.Context(), s.db, &f); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := encoding.Respond(w, r, &f, http.StatusCreated); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
