package httphandler

import (
	"context"
	"net/http"
	"upload/pkg/file"
	"upload/util/blob"
	"upload/util/encoding"

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

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		errChan := make(chan error)
		go func() { errChan <- file.DeleteByID(ctx, s.db, uuid) }()
		go func() { errChan <- blob.Delete(ctx, "company", uuid.String()) }()
		if err := <-errChan; err != nil {
			encoding.ErrorRespond(w, r, http.StatusInternalServerError, file.ErrInternal)
			return
		}
		if err := <-errChan; err != nil {
			encoding.ErrorRespond(w, r, http.StatusInternalServerError, file.ErrInternal)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
