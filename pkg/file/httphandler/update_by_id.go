package httphandler

// import (
// 	"net/http"
// 	"upload/pkg/file"
// 	"upload/util/encoding"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/google/uuid"
// )

// func (s *fileServer) handleFileUpdate() http.HandlerFunc {
// 	type request struct {
// 		UploaderID  string `json:"uploaderId" validate:"uuid"`
// 		CompanyID   string `json:"companyId" validate:"uuid"`
// 		Description string `json:"description"`
// 	}

// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var req request

// 		if err := encoding.Decode(w, r, &req); err != nil {
// 			encoding.ErrorRespond(w, r, http.StatusBadRequest, err)
// 			return
// 		}

// 		if err := s.validator.Struct(req); err != nil {
// 			encoding.ErrorRespond(w, r, http.StatusBadRequest, err)
// 			return
// 		}

// 		var f file.File
// 		f.UploaderID = req.UploaderID
// 		f.CompanyID = req.CompanyID
// 		f.Description = req.Description

// 		id := chi.URLParam(r, "fileID")
// 		uuid, err := uuid.Parse(id)
// 		f.ID = uuid
// 		if err != nil {
// 			encoding.ErrorRespond(w, r, http.StatusBadRequest, err)
// 			return
// 		}

// 		if err := file.UpdateByID(r.Context(), s.db, uuid, &f); err != nil {
// 			encoding.ErrorRespond(w, r, http.StatusInternalServerError, file.ErrInternal)
// 			return
// 		}

// 		resp := DataResponse{&f}
// 		if err := encoding.Respond(w, r, resp, http.StatusOK); err != nil {
// 			encoding.ErrorRespond(w, r, http.StatusInternalServerError, file.ErrInternal)
// 			return
// 		}
// 	}
// }
