package file

import (
	"time"
	"upload/util/enum"

	"github.com/google/uuid"
)

type File struct {
	ID              uuid.UUID         `json:"id"`
	UploaderID      string            `json:"uploaderId"`
	CompanyID       string            `json:"companyId"`
	Name            string            `json:"name"`
	Extension       string            `json:"extension"`
	ContentType     string            `json:"contentType"`
	Size            uint              `json:"size"`
	StorageLocation string            `json:"storageLocation"`
	TimesRequested  uint              `json:"timesRequested"`
	Status          enum.UploadStatus `json:"status"`
	Title           string            `json:"title"`
	Description     string            `json:"description"`
	SubmittedAt     time.Time         `json:"submittedAt"`
	UpdatedAt       time.Time         `json:"updatedAt"`
	UploadedAt      time.Time         `json:"uploadedAt"`
}
