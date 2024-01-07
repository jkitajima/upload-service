package mongo

import (
	"context"
	"time"

	"upload/pkg/file"
	"upload/util/enum"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (db *FileCollection) Insert(ctx context.Context, f *file.File) error {
	f.ID = uuid.New()

	now := time.Now()
	f.UpdatedAt = now
	f.UploadedAt = now

	binID := primitive.Binary{
		Subtype: bson.TypeBinaryUUID,
		Data:    []byte(f.ID[:]),
	}

	doc := struct {
		ID              primitive.Binary  `bson:"_id"`
		UploaderID      string            `bson:"uploaderId"`
		CompanyID       string            `bson:"companyId"`
		Name            string            `bson:"name"`
		Extension       string            `bson:"extension"`
		ContentType     string            `bson:"contentType"`
		Size            uint              `bson:"size"`
		StorageLocation string            `bson:"storageLocation"`
		TimesRequested  uint              `bson:"timesRequested"`
		Status          enum.UploadStatus `bson:"status"`
		Title           string            `bson:"title"`
		Description     string            `bson:"description"`
		SubmittedAt     time.Time         `bson:"submittedAt"`
		UpdatedAt       time.Time         `bson:"updatedAt"`
		UploadedAt      time.Time         `bson:"uploadedAt"`
	}{
		binID,
		f.UploaderID,
		f.CompanyID,
		f.Name,
		f.Extension,
		f.ContentType,
		f.Size,
		f.StorageLocation,
		f.TimesRequested,
		f.Status,
		f.Title,
		f.Description,
		f.SubmittedAt,
		f.UpdatedAt,
		f.UploadedAt,
	}

	_, err := db.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	return nil
}
