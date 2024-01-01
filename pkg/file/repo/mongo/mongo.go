package mongo

import "go.mongodb.org/mongo-driver/mongo"

type FileCollection struct {
	*mongo.Collection
}

func NewRepo(db *mongo.Collection) *FileCollection {
	return &FileCollection{db}
}
