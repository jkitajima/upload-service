package httphandler

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"

	"upload/pkg/file"
	"upload/util/blob"
	"upload/util/encoding"

	blobOpts "gocloud.dev/blob"
)

func (s *fileServer) handleFileCreate() http.HandlerFunc {
	const maxPermittedRequestSize = 4 << 20

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

		// insert into db and blob stg
		ctx := r.Context()
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

			err = blob.Upload(ctx, "company", uploadedFile.Filename, buf, &opts)
			if err != nil {
				errChan <- err
				return
			}

			errChan <- nil
		}()

		go func() {
			if err := file.Create(ctx, s.db, &f); err != nil {
				errChan <- err
				return
			}

			errChan <- nil
		}()

		// check for blob stg and db err
		if err := <-errChan; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := <-errChan; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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
		"name":       false,
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
		case "name":
			if len(v) > 1 {
				return errors.New(`must have only one "name" form field`)
			}
			f.Name = v[0]
			requiredFields[k] = true
			fieldsCounter--
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
