package mongo

import (
	"context"
	"errors"
	"fmt"
	"upload/shared/zombiekiller"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (db *FileCollection) KillZombie(key fmt.Stringer) error {
	id, err := uuid.Parse(key.String())
	if err != nil {
		return errors.New("mongo: kill zombie: key must be a valid UUID")
	}

	binID := primitive.Binary{
		Subtype: bson.TypeBinaryUUID,
		Data:    []byte(id[:]),
	}

	filter := bson.D{{
		Key:   "_id",
		Value: binID,
	}}

	result, err := db.DeleteOne(context.TODO(), filter)
	if err != nil {
		return zombiekiller.ErrInternal
	}

	if result.DeletedCount == 0 {
		return zombiekiller.ErrNotFound
	}

	return nil
}
