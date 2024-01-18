package mongo

import (
	"context"
	"log"

	"upload/pkg/file"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (db *FileCollection) FindByID(ctx context.Context, id uuid.UUID) (*file.File, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var f file.File
	f.ID = id

	binID := primitive.Binary{
		Subtype: bson.TypeBinaryUUID,
		Data:    []byte(id[:]),
	}

	incChan := make(chan error)
	go func() {
		_, err := db.Collection.UpdateByID(ctx, binID, bson.D{{Key: "$inc", Value: bson.D{{Key: "timesRequested", Value: 1}}}})
		incChan <- err
	}()

	findChan := make(chan error)
	go func() {
		findChan <- db.FindOne(
			ctx, bson.D{{
				Key:   "_id",
				Value: binID,
			}},
		).Decode(&f)
	}()

	var errored bool
	var errResponse error

	select {
	case err := <-incChan:
		if err != nil {
			cancel()
			log.Printf("file: repo: mongo: find_by_id: inc: %v\n", err)
		}
	case err := <-findChan:
		if err != nil {
			cancel()
			log.Printf("file: repo: mongo: find_by_id: %v\n", err)
			errored = true
			errResponse = file.ErrInternal
			if err == mongo.ErrNoDocuments {
				errResponse = file.ErrNotFoundByID
			}
		}
	}

	select {
	case err := <-incChan:
		if err != nil {
			log.Printf("file: repo: mongo: find_by_id: inc: %v\n", err)
		}

		f.TimesRequested++
	case err := <-findChan:
		if err != nil {
			log.Printf("file: repo: mongo: find_by_id: %v\n", err)
			errored = true
			errResponse = err
			if err == mongo.ErrNoDocuments {
				errResponse = file.ErrNotFoundByID
			}
		}
	}

	if errored {
		return nil, errResponse
	}
	return &f, nil
}
