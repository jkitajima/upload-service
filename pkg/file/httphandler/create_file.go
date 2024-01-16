package httphandler

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"upload/pkg/file"
	"upload/shared/encoding"

	"github.com/google/uuid"
)

func (s *fileServer) handleFileCreate() http.HandlerFunc {
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

	const maxPermittedRequestSize = 4 << 20 // 4 MiB
	var now = time.Now()

	return func(w http.ResponseWriter, r *http.Request) {
		// validate max request size
		requestSize := r.ContentLength
		if requestSize > maxPermittedRequestSize {
			encoding.ErrorRespond(w, r, http.StatusBadRequest, fmt.Errorf("maximum permitted request size is %d bytes", maxPermittedRequestSize))
			return
		}

		// validate if request is "multipart/form-data"
		if err := r.ParseMultipartForm(maxPermittedRequestSize); err != nil {
			encoding.ErrorRespond(w, r, http.StatusBadRequest, errors.New("failed to parse multipart/form-data input"))
			return
		}

		// validate and fetch uploaded file
		uploadedFile, err := fetchFormFile(r)
		if err != nil {
			encoding.ErrorRespond(w, r, http.StatusBadRequest, err)
			return
		}

		// validate form field values and build a File struct
		var f file.File
		if err := validateFormValues(r, &f); err != nil {
			encoding.ErrorRespond(w, r, http.StatusBadRequest, err)
			return
		}

		// now that input is validated
		// add server-side file attributes
		f.ID = uuid.New()
		f.Name = uploadedFile.Filename
		f.Extension = filepath.Ext(uploadedFile.Filename)

		const blobstgContainer = "company"
		f.ContentType = uploadedFile.Header.Get("Content-Type")
		f.Size = uint(uploadedFile.Size)
		f.SubmittedAt = now
		f.StorageLocation = fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", os.Getenv("AZURE_STORAGE_ACCOUNT"), blobstgContainer, f.ID.String())
		uploadedFile.Filename = f.ID.String()

		// open file to be uploaded
		openedFile, _ := uploadedFile.Open()
		defer openedFile.Close()

		service := file.Service{Repo: s.db, Blob: s.blobstg}
		serviceRequest := file.CreateRequest{
			Metadata: &f,
			Rawdata:  openedFile,
			Bucket:   blobstgContainer,
		}

		serviceResponse, err := service.Create(r.Context(), serviceRequest)
		if err != nil {
			switch err {
			default:
				encoding.ErrorRespond(w, r, http.StatusInternalServerError, err)
				return
			}
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
		if err := encoding.Respond(w, r, &dataresp, http.StatusCreated); err != nil {
			encoding.ErrorRespond(w, r, http.StatusInternalServerError, file.ErrInternal)
			return
		}
	}
}

func fetchFormFile(r *http.Request) (*multipart.FileHeader, error) {
	filesUploaded := len(r.MultipartForm.File)
	if filesUploaded == 0 {
		return nil, errors.New("at least one form field should contain file input")
	} else if filesUploaded > 1 {
		return nil, errors.New("only one form field can contain file input")
	}

	fileForm, ok := r.MultipartForm.File["file"]
	if !ok {
		return nil, errors.New(`multipart must have a form field named "file"`)
	} else if len(fileForm) == 0 {
		return nil, errors.New(`form field "file" must have a file input`)
	} else if len(fileForm) > 1 {
		return nil, errors.New(`form field "file" can have only a single file input`)
	}

	return fileForm[0], nil
}

func validateFormValues(r *http.Request, f *file.File) error {
	var requiredFields = map[string]bool{
		"uploaderId": false,
		"companyId":  false,
	}
	var fieldsCounter int = len(requiredFields)

	for k, v := range r.MultipartForm.Value {
		switch k {
		case "uploaderId":
			if len(v) > 1 {
				return errors.New(`must have only one "uploaderId" form field`)
			}
			f.UploaderID = v[0]
			requiredFields[k] = true
			fieldsCounter--
		case "companyId":
			if len(v) > 1 {
				return errors.New(`must have only one "companyId" form field`)
			}
			f.CompanyID = v[0]
			requiredFields[k] = true
			fieldsCounter--
		case "description":
			if len(v) > 1 {
				return errors.New(`must have only one "description" form field`)
			}
			f.Description = v[0]
		}
	}

	if fieldsCounter > 0 {
		missingFields := "missing fields: "
		for k, v := range requiredFields {
			if !v {
				missingFields += k + ", "
			}
		}
		missingFields = missingFields[:len(missingFields)-2]

		return errors.New(missingFields)
	}

	return nil
}
