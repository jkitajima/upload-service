package httphandler

import (
	"context"
	"net/http"
	"upload/pkg/file"
	"upload/util/blob"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (s *fileServer) handleFileDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "fileID")
		uuid, err := uuid.Parse(id)
		if err != nil {
			http.Error(w, "failed to parse uuid", http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		errChan := make(chan error)
		go func() { errChan <- file.Delete(ctx, s.db, uuid) }()
		go func() { errChan <- blob.Delete(ctx, "company", uuid.String()) }()
		if err := <-errChan; err != nil {
			http.Error(w, "failed to delete requested file", http.StatusInternalServerError)
			return
		}
		if err := <-errChan; err != nil {
			http.Error(w, "failed to delete requested file", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
