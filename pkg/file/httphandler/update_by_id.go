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
			resp := NewErrorsResponse(&ErrorObject{http.StatusBadRequest, "Invalid Input", "Failed to decode client input."})
			encoding.Respond(w, r, resp, http.StatusBadRequest)
			return
		}

		id := chi.URLParam(r, "fileID")
		uuid, err := uuid.Parse(id)
		f.ID = uuid
		if err != nil {
			resp := NewErrorsResponse(&ErrorObject{http.StatusBadRequest, "Invalid Input", "Failed to parse UUID."})
			encoding.Respond(w, r, resp, http.StatusBadRequest)
			return
		}

		if err := file.Update(r.Context(), s.db, uuid, &f); err != nil {
			resp := NewErrorsResponse(&ErrorObject{http.StatusInternalServerError, "Internal Server Error", "Server encountered an unexpected condition that prevented it from fulfilling the request."})
			encoding.Respond(w, r, resp, http.StatusInternalServerError)
			return
		}

		resp := DataResponse{&f}
		if err := encoding.Respond(w, r, resp, http.StatusOK); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
