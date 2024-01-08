package mongo

import (
	"context"

	"upload/pkg/file"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (db *FileCollection) FindByID(ctx context.Context, id uuid.UUID) (*file.File, error) {
	var f file.File
	f.ID = id

	binID := primitive.Binary{
		Subtype: bson.TypeBinaryUUID,
		Data:    []byte(id[:]),
	}

	incChan := make(chan error)
	findChan := make(chan error)
	go func() {
		_, err := db.UpdateByID(ctx, binID, bson.D{{Key: "$inc", Value: bson.D{{Key: "timesRequested", Value: 1}}}})
		incChan <- err
	}()

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
			return nil, err
		}
	case err := <-findChan:
		if err != nil {
			return nil, err
		}
	}

	select {
	case err := <-incChan:
		if err != nil {
			return nil, err
		}

		f.TimesRequested++
	case err := <-findChan:
		if err != nil {
			return nil, err
		}
	}

	return &f, nil
}
