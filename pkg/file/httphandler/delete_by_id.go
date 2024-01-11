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
			resp := NewErrorsResponse(&ErrorObject{http.StatusBadRequest, "Invalid Input", "Failed to parse UUID."})
			encoding.Respond(w, r, resp, http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		errChan := make(chan error)
		go func() { errChan <- file.Delete(ctx, s.db, uuid) }()
		go func() { errChan <- blob.Delete(ctx, "company", uuid.String()) }()
		if err := <-errChan; err != nil {
			resp := NewErrorsResponse(&ErrorObject{http.StatusInternalServerError, "Internal Server Error", "Server encountered an unexpected condition that prevented it from fulfilling the request."})
			encoding.Respond(w, r, resp, http.StatusInternalServerError)
			return
		}
		if err := <-errChan; err != nil {
			resp := NewErrorsResponse(&ErrorObject{http.StatusInternalServerError, "Internal Server Error", "Server encountered an unexpected condition that prevented it from fulfilling the request."})
			encoding.Respond(w, r, resp, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
