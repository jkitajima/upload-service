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

	select {
	case err := <-incChan:
		if err != nil {
			log.Println(err)
			return nil, file.ErrInternal
		}
	case err := <-findChan:
		if err != nil {
			log.Println(err)

			if err == mongo.ErrNoDocuments {
				return nil, file.ErrFileNotFoundByID
			}

			return nil, file.ErrInternal
		}
	}

	select {
	case err := <-incChan:
		if err != nil {
			log.Println(err)
			return nil, file.ErrInternal
		}

		f.TimesRequested++
	case err := <-findChan:
		if err != nil {
			log.Println(err)

			if err == mongo.ErrNoDocuments {
				return nil, file.ErrFileNotFoundByID
			}

			return nil, file.ErrInternal
		}
	}

	return &f, nil
}
