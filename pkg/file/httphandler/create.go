package httphandler

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"upload/pkg/file"
	"upload/util/blob"
	"upload/util/encoding"

	"github.com/google/uuid"
	blobOpts "gocloud.dev/blob"
)

func (s *fileServer) handleFileCreate() http.HandlerFunc {
	const maxPermittedRequestSize = 4 << 20
	var now = time.Now()

	return func(w http.ResponseWriter, r *http.Request) {
		// validate max request size
		requestSize := r.ContentLength
		if requestSize > maxPermittedRequestSize {
			http.Error(w, fmt.Sprintf("max permitted request size is: %d bytes", maxPermittedRequestSize), http.StatusBadRequest)
			return
		}

		// validate if request is "multipart/form-data"
		if err := r.ParseMultipartForm(maxPermittedRequestSize); err != nil {
			http.Error(w, http.ErrNotMultipart.Error(), http.StatusBadRequest)
			return
		}

		// validate and fetch uploaded file
		uploadedFile, err := fetchFormFile(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// validate form field values and build a File struct
		var f file.File
		if err := validateFormValues(r, &f); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// server-side file attributes
		f.ID = uuid.New()
		f.Name = uploadedFile.Filename
		f.Extension = filepath.Ext(uploadedFile.Filename)

		const blobstgContainer = "company"
		f.ContentType = uploadedFile.Header.Get("Content-Type")
		f.Size = uint(uploadedFile.Size)
		f.SubmittedAt = now
		f.StorageLocation = fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", os.Getenv("AZURE_STORAGE_ACCOUNT"), blobstgContainer, f.ID.String())
		uploadedFile.Filename = f.ID.String()

		// insert into db and blob stg
		ctx := r.Context()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		blobstgChan := uploadToBlobStorage(ctx, uploadedFile, blobstgContainer)
		dbChan := insertIntoDB(ctx, s, &f)

	RangeChannels:
		for i := 0; i < 2; i++ {
			select {
			case err := <-blobstgChan:
				if err != nil {
					cancel()
					http.Error(w, err.Error(), http.StatusInternalServerError)
					break RangeChannels
				}
			case err := <-dbChan:
				if err != nil {
					cancel()
					http.Error(w, err.Error(), http.StatusInternalServerError)
					break RangeChannels
				}
			}
		}

		if err := encoding.Respond(w, r, &f, http.StatusCreated); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func fetchFormFile(r *http.Request) (*multipart.FileHeader, error) {
	filesUploaded := len(r.MultipartForm.File)
	if filesUploaded == 0 {
		return nil, errors.New("at least one form field should contain files")
	} else if filesUploaded > 1 {
		return nil, errors.New("only one form field should contain files")
	}

	fileForm, ok := r.MultipartForm.File["file"]
	if !ok {
		return nil, errors.New(`multipart must have a form field name "file"`)
	} else if len(fileForm) == 0 {
		return nil, errors.New(`"file" form field must have a file`)
	} else if len(fileForm) > 1 {
		return nil, errors.New(`"file" form field must only have a single file`)
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

func blobstg(ctx context.Context, uploadedFile *multipart.FileHeader, bucket string) <-chan error {
	errChan := make(chan error)

	go func() {
		openedFile, err := uploadedFile.Open()
		if err != nil {
			errChan <- err
			return
		}
		defer openedFile.Close()

		buf := make([]byte, uploadedFile.Size)
		_, err = openedFile.Read(buf)
		if err != nil {
			errChan <- err
			return
		}

		var opts blobOpts.WriterOptions
		header := uploadedFile.Header
		contentType := header.Get("Content-Type")
		if contentType != "" {
			opts.ContentType = contentType
		}

		err = blob.Upload(ctx, bucket, uploadedFile.Filename, buf, &opts)
		if err != nil {
			errChan <- err
			return
		}

		errChan <- nil
	}()

	return errChan
}

func uploadToBlobStorage(ctx context.Context, uploadedFile *multipart.FileHeader, bucket string) <-chan error {
	errChan := make(chan error)

	go func() {
		blobstgChan := blobstg(ctx, uploadedFile, bucket)

		select {
		case <-ctx.Done():
			errChan <- ctx.Err()
		case err := <-blobstgChan:
			errChan <- err
		}
	}()

	return errChan
}

func db(ctx context.Context, s *fileServer, f *file.File) <-chan error {
	errChan := make(chan error)

	go func() {
		if err := file.Create(ctx, s.db, f); err != nil {
			errChan <- err
			return
		}

		errChan <- nil
	}()

	return errChan
}

func insertIntoDB(ctx context.Context, s *fileServer, f *file.File) <-chan error {
	errChan := make(chan error)

	go func() {
		dbChan := db(ctx, s, f)

		select {
		case <-ctx.Done():
			errChan <- ctx.Err()
		case err := <-dbChan:
			errChan <- err
		}
	}()

	return errChan
}
