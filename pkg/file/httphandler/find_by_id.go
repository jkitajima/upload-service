package httphandler

import (
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
			resp := NewErrorsResponse(&ErrorObject{http.StatusBadRequest, "Invalid Path Parameter", "File ID is in a invalid format."})
			encoding.Respond(w, r, resp, http.StatusBadRequest)
			return
		}

		f, err := file.FindByID(r.Context(), s.db, uuid)
		if err != nil {
			http.Error(w, "could not find any file with provided id", http.StatusInternalServerError)
			return
		}

		resp := DataResponse{f}
		if err := encoding.Respond(w, r, resp, http.StatusOK); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
