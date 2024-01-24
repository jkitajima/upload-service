package mongo

import (
	"context"
	"log"
	"time"

	"upload/pkg/file"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (db *FileCollection) Insert(ctx context.Context, f *file.File) error {
	now := time.Now()
	f.UpdatedAt = now
	f.UploadedAt = now

	binID := primitive.Binary{
		Subtype: bson.TypeBinaryUUID,
		Data:    []byte(f.ID[:]),
	}

	doc := struct {
		ID              primitive.Binary `bson:"_id"`
		UploaderID      string           `bson:"uploaderId"`
		CompanyID       string           `bson:"companyId"`
		Name            string           `bson:"name"`
		Extension       string           `bson:"extension"`
		ContentType     string           `bson:"contentType"`
		Size            uint             `bson:"size"`
		StorageLocation string           `bson:"storageLocation"`
		TimesRequested  uint             `bson:"timesRequested"`
		Description     string           `bson:"description"`
		SubmittedAt     time.Time        `bson:"submittedAt"`
		UpdatedAt       time.Time        `bson:"updatedAt"`
		UploadedAt      time.Time        `bson:"uploadedAt"`
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
		f.Description,
		f.SubmittedAt,
		f.UpdatedAt,
		f.UploadedAt,
	}

	_, err := db.InsertOne(ctx, doc)
	if err != nil {
		log.Printf("file: repo: mongo: insert: %v\n", err)
		if err := mongo.IsDuplicateKeyError(err); err {
			return file.ErrInsertDuplicatedKey
		}
		return file.ErrInternal
	}

	return nil
}
