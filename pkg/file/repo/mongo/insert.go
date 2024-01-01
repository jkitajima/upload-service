package mongo

import (
	"context"

	"upload/pkg/file"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (db *FileCollection) Insert(ctx context.Context, f *file.File) error {
	f.ID = uuid.New()

	binID := primitive.Binary{
		Subtype: bson.TypeBinaryUUID,
		Data:    []byte(f.ID[:]),
	}

	doc := struct {
		ID   primitive.Binary `bson:"_id"`
		Name string
	}{
		binID,
		f.Name,
	}

	_, err := db.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	return nil
}
