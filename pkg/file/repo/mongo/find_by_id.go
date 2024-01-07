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

	_, err := db.UpdateByID(ctx, binID, bson.D{{Key: "$inc", Value: bson.D{{Key: "timesRequested", Value: 1}}}})
	if err != nil {
		return nil, err
	}

	err = db.FindOne(
		ctx, bson.D{{
			Key:   "_id",
			Value: binID,
		}},
	).Decode(&f)

	if err != nil {
		log.Println("failed to decode result")

		if err == mongo.ErrNoDocuments {
			log.Println("given filter did not match any docs on collection")
		}

		return nil, err
	}

	return &f, nil
}
