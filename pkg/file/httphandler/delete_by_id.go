package httphandler

// import (
// 	"context"
// 	"log"
// 	"net/http"
// 	"upload/pkg/file"
// 	"upload/util/blob"
// 	"upload/util/encoding"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/google/uuid"
// )

// func (s *fileServer) handleFileDelete() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		id := chi.URLParam(r, "fileID")
// 		uuid, err := uuid.Parse(id)
// 		if err != nil {
// 			encoding.ErrorRespond(w, r, http.StatusBadRequest, err)
// 			return
// 		}

// 		ctx, cancel := context.WithCancel(r.Context())
// 		defer cancel()

// 		dbChan := make(chan error)
// 		go func() { dbChan <- file.DeleteByID(ctx, s.db, uuid) }()

// 		blobChan := make(chan error)
// 		go func() { blobChan <- blob.Delete(ctx, "company", uuid.String()) }()

// 		for i := 0; i < 2; i++ {
// 			select {
// 			case err := <-dbChan:
// 				if err != nil {
// 					log.Println(err)
// 					cancel()
// 					if err == file.ErrFileNotFoundByID {
// 						encoding.ErrorRespond(w, r, http.StatusBadRequest, err)
// 						return
// 					}
// 					encoding.ErrorRespond(w, r, http.StatusInternalServerError, file.ErrInternal)
// 					return
// 				}
// 			case err := <-blobChan:
// 				if err != nil {
// 					log.Println(err)
// 					cancel()
// 					if err == blob.ErrNotFound {
// 						encoding.ErrorRespond(w, r, http.StatusBadRequest, file.ErrFileNotFoundByID)
// 						return
// 					}
// 					encoding.ErrorRespond(w, r, http.StatusInternalServerError, file.ErrInternal)
// 					return
// 				}
// 			}
// 		}

// 		w.WriteHeader(http.StatusNoContent)
// 	}
// }
