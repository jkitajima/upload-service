package httphandler

import (
	"errors"
	"log"
	"net/http"
	"time"

	"upload/pkg/file"
	"upload/shared/encoding"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (s *fileServer) handleFileFindByID() http.HandlerFunc {
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
		id := chi.URLParam(r, "fileID")

		uuid, err := uuid.Parse(id)
		if err != nil {
			encoding.ErrorRespond(w, r, http.StatusBadRequest, err)
			return
		}

		service := file.Service{Repo: s.db, Blob: s.blobstg, Thrash: s.thrash}
		serviceRequest := file.FindByIDRequest{ID: uuid}
		serviceResponse, err := service.FindByID(r.Context(), serviceRequest)
		if err != nil {
			switch err {
			case file.ErrNotFoundByID:
				encoding.ErrorRespond(w, r, http.StatusNotFound, err)
			default:
				encoding.ErrorRespond(w, r, http.StatusInternalServerError, err)
			}
			return
		}

		params := r.URL.Query()
		if len(params) > 0 {
			switch params.Get("redirect") {
			case "storageLocation":
				redirect := serviceResponse.Metadata.StorageLocation
				log.Printf("redirecting client to %q", redirect)
				http.Redirect(w, r, redirect, http.StatusSeeOther)
			default:
				encoding.ErrorRespond(w, r, http.StatusBadRequest, errors.New("value of `redirect` parameter must be a valid file attribute: [`storageLocation`]"))
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

		dataresp := DataResponse{&resp}
		if err := encoding.Respond(w, r, dataresp, http.StatusOK); err != nil {
			encoding.ErrorRespond(w, r, http.StatusInternalServerError, file.ErrInternal)
			return
		}
	}
}
