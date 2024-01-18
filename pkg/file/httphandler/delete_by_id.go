package httphandler

import (
	"net/http"

	"upload/pkg/file"
	"upload/shared/encoding"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (s *fileServer) handleFileDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "fileID")
		uuid, err := uuid.Parse(id)
		if err != nil {
			encoding.ErrorRespond(w, r, http.StatusBadRequest, err)
			return
		}

		service := file.Service{Repo: s.db, Blob: s.blobstg, Thrash: s.thrash}
		serviceRequest := file.DeleteByIDRequest{ID: uuid, Bucket: "company"}
		if err := service.DeleteByID(r.Context(), serviceRequest); err != nil {
			switch err {
			case file.ErrNotFoundByID:
				encoding.ErrorRespond(w, r, http.StatusNotFound, err)
			default:
				encoding.ErrorRespond(w, r, http.StatusInternalServerError, err)
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
