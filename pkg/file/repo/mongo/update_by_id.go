package mongo

import (
	"context"
	"log"
	"time"

	"upload/pkg/file"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (db *FileCollection) UpdateByID(ctx context.Context, id uuid.UUID, f *file.File) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	binID := primitive.Binary{
		Subtype: bson.TypeBinaryUUID,
		Data:    []byte(id[:]),
	}

	c := make(Changes)
	updatedAt := time.Now()
	update := bson.D{{
		Key:   "$set",
		Value: c.retrieveNonNullableElements(f, updatedAt),
	}}

	updateChan := make(chan error)
	go func() {
		_, err := db.Collection.UpdateByID(ctx, binID, update)
		updateChan <- err
	}()

	findChan := make(chan error)
	go func() {
		findChan <- db.FindOne(
			ctx, bson.D{{
				Key:   "_id",
				Value: binID,
			}},
		).Decode(f)
	}()

	var errored bool
	var errResponse error

	select {
	case err := <-updateChan:
		if err != nil {
			cancel()
			log.Printf("file: repo: mongo: update_by_id: update: %v\n", err)
			errored = true
			errResponse = file.ErrInternal
		}
	case err := <-findChan:
		if err != nil {
			cancel()
			log.Printf("file: repo: mongo: update_by_id: find: %v\n", err)
			errored = true
			errResponse = file.ErrInternal
			if err == mongo.ErrNoDocuments {
				errResponse = file.ErrNotFoundByID
			}
		}
	}

	select {
	case err := <-updateChan:
		if err != nil {
			log.Printf("file: repo: mongo: update_by_id: update: %v\n", err)
			errored = true
			errResponse = file.ErrInternal
		}

		if err == nil && !errored {
			if v, ok := c[uploaderID]; ok {
				f.UploaderID = v.(string)
			}

			if v, ok := c[companyID]; ok {
				f.CompanyID = v.(string)
			}

			if v, ok := c[description]; ok {
				f.Description = v.(string)
			}

			f.UpdatedAt = updatedAt
		}
	case err := <-findChan:
		if err != nil {
			log.Printf("file: repo: mongo: update_by_id: find: %v\n", err)
			errored = true
			errResponse = file.ErrInternal
			if err == mongo.ErrNoDocuments {
				errResponse = file.ErrNotFoundByID
			}
		}
	}

	if errored {
		return errResponse
	}
	return nil
}

const (
	uploaderID  = "uploaderId"
	companyID   = "companyId"
	description = "description"
	updatedAt   = "updatedAt"
)

type Changes map[string]any

func (c Changes) retrieveNonNullableElements(f *file.File, t time.Time) []primitive.E {
	elements := make([]primitive.E, 0, 3)

	if f.UploaderID != "" {
		elements = append(elements, primitive.E{Key: uploaderID, Value: f.UploaderID})
		c[uploaderID] = f.UploaderID
	}

	if f.CompanyID != "" {
		elements = append(elements, primitive.E{Key: companyID, Value: f.CompanyID})
		c[companyID] = f.CompanyID
	}

	if f.Description != "" {
		elements = append(elements, primitive.E{Key: description, Value: f.Description})
		c[description] = f.Description
	}

	elements = append(elements, primitive.E{Key: updatedAt, Value: t})
	return elements
}
