package mongo

import (
	"context"

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

	result, err := db.DeleteOne(ctx, filter)
	if err != nil {
		return file.ErrInternal
	}
	if result.DeletedCount == 0 {
		return file.ErrNotFoundByID
	}

	return nil
}
