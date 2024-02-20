package httphandler

import (
	"net/http"
	"time"

	"upload/pkg/file"
	"upload/shared/encoding"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (s *fileServer) handleFileUpdate() http.HandlerFunc {
	type request struct {
		UploaderID  string `json:"uploaderId"`
		CompanyID   string `json:"companyId"`
		Description string `json:"description"`
	}

	type response struct {
		ID              uuid.UUID `json:"id"`
		UploaderID      string    `json:"uploaderId"`
		CompanyID       string    `json:"companyId"`
		Name            string    `json:"name"`
		Extension       string    `json:"extension"`
		ContentType     string    `json:"contentType"`
		Size            uint      `json:"size"`
		StorageLocation string    `json:"storageLocation"`
		TimesRequested  uint      `json:"timesRequested"`
		Description     string    `json:"description"`
		SubmittedAt     time.Time `json:"submittedAt"`
		UpdatedAt       time.Time `json:"updatedAt"`
		UploadedAt      time.Time `json:"uploadedAt"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request

		if err := encoding.Decode(w, r, &req); err != nil {
			encoding.ErrorRespond(w, r, http.StatusBadRequest, err)
			return
		}

		var f file.File
		f.UploaderID = req.UploaderID
		f.CompanyID = req.CompanyID
		f.Description = req.Description

		id := chi.URLParam(r, "fileID")
		uuid, err := uuid.Parse(id)
		f.ID = uuid
		if err != nil {
			encoding.ErrorRespond(w, r, http.StatusBadRequest, err)
			return
		}

		service := file.Service{Repo: s.db, Blob: s.blobstg, Thrash: s.thrash}
		serviceRequest := file.UpdateByIDRequest{ID: uuid, Metadata: &f}
		serviceResponse, err := service.UpdateByID(r.Context(), serviceRequest)
		if err != nil {
			switch err {
			case file.ErrNotFoundByID:
				encoding.ErrorRespond(w, r, http.StatusNotFound, err)
			default:
				encoding.ErrorRespond(w, r, http.StatusInternalServerError, err)
			}
			return
		}

		meta := serviceResponse.Metadata
		resp := response{
			ID:              meta.ID,
			UploaderID:      meta.UploaderID,
			CompanyID:       meta.CompanyID,
			Name:            meta.Name,
			Extension:       meta.Extension,
			ContentType:     meta.ContentType,
			Size:            meta.Size,
			StorageLocation: meta.StorageLocation,
			TimesRequested:  meta.TimesRequested,
			Description:     meta.Description,
			SubmittedAt:     meta.SubmittedAt,
			UpdatedAt:       meta.UpdatedAt,
			UploadedAt:      meta.UploadedAt,
		}

		dataresp := encoding.DataResponse{Data: &resp}
		if err := encoding.Respond(w, r, &dataresp, http.StatusOK); err != nil {
			encoding.ErrorRespond(w, r, http.StatusInternalServerError, file.ErrInternal)
			return
		}
	}
}
