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
			resp := NewErrorsResponse(&ErrorObject{http.StatusNotFound, "Targeted Resource Not Found", "Could not find any file with provided id."})
			encoding.Respond(w, r, resp, http.StatusNotFound)
			return
		}

		params := r.URL.Query()
		if len(params) > 0 {
			switch params.Get("redirect") {
			case "storageLocation":
				http.Redirect(w, r, f.StorageLocation, http.StatusSeeOther)
			default:
				resp := NewErrorsResponse(&ErrorObject{http.StatusBadRequest, "Invalid Query Parameter", "Value of `redirect` parameter must be a valid file attribute: [`storageLocation`]."})
				encoding.Respond(w, r, resp, http.StatusBadRequest)
				return
			}
		}

		resp := DataResponse{f}
		if err := encoding.Respond(w, r, resp, http.StatusOK); err != nil {
			resp := NewErrorsResponse(&ErrorObject{http.StatusInternalServerError, "Internal Server Error", "Server encountered an unexpected condition that prevented it from fulfilling the request."})
			encoding.Respond(w, r, resp, http.StatusInternalServerError)
			return
		}
	}
}
