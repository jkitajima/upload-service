package mongo

import (
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (db *FileCollection) Delete(ctx context.Context, id uuid.UUID) error {
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
		return err
	}

	return nil
}
