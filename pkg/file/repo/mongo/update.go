package mongo

import (
	"context"

	"upload/pkg/file"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (db *FileCollection) Update(ctx context.Context, id uuid.UUID, f *file.File) error {
	binID := primitive.Binary{
		Subtype: bson.TypeBinaryUUID,
		Data:    []byte(id[:]),
	}

	update := bson.D{{
		Key:   "$set",
		Value: bson.D{{Key: "name", Value: f.Name}},
	}}

	_, err := db.UpdateByID(ctx, binID, update)
	if err != nil {
		return err
	}

	return nil
}
