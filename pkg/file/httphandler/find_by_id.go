package httphandler

// import (
// 	"errors"
// 	"net/http"

// 	"upload/pkg/file"
// 	"upload/util/encoding"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/google/uuid"
// )

// func (s *fileServer) handleFileFindByID() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		id := chi.URLParam(r, "fileID")

// 		uuid, err := uuid.Parse(id)
// 		if err != nil {
// 			encoding.ErrorRespond(w, r, http.StatusBadRequest, err)
// 			return
// 		}

// 		f, err := file.FindByID(r.Context(), s.db, uuid)
// 		if err != nil {
// 			switch err {
// 			case file.ErrFileNotFoundByID:
// 				encoding.ErrorRespond(w, r, http.StatusNotFound, err)
// 			}

// 			return
// 		}

// 		params := r.URL.Query()
// 		if len(params) > 0 {
// 			switch params.Get("redirect") {
// 			case "storageLocation":
// 				http.Redirect(w, r, f.StorageLocation, http.StatusSeeOther)
// 			default:
// 				encoding.ErrorRespond(w, r, http.StatusBadRequest, errors.New("value of `redirect` parameter must be a valid file attribute: [`storageLocation`]"))
// 				return
// 			}
// 		}

// 		resp := encoding.DataResponse{Data: f}
// 		if err := encoding.Respond(w, r, resp, http.StatusOK); err != nil {
// 			encoding.ErrorRespond(w, r, http.StatusInternalServerError, file.ErrInternal)
// 			return
// 		}
// 	}
// }
