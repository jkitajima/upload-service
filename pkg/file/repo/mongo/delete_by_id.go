package mongo

import (
	"context"
	"log"
	"upload/pkg/file"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (db *FileCollection) DeleteByID(ctx context.Context, id uuid.UUID) error {
	binID := primitive.Binary{
		Subtype: bson.TypeBinaryUUID,
		Data:    []byte(id[:]),
	}

	filter := bson.D{{
		Key:   "_id",
		Value: binID,
	}}

	_, err := db.DeleteOne(ctx, filter)
	if err != nil {
		log.Println(err)
		return file.ErrInternal
	}

	return nil
}
