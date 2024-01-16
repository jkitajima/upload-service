package mongo

// import (
// 	"context"
// 	"log"
// 	"upload/pkg/file"

// 	"github.com/google/uuid"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// func (db *FileCollection) DeleteByID(ctx context.Context, id uuid.UUID) error {
// 	binID := primitive.Binary{
// 		Subtype: bson.TypeBinaryUUID,
// 		Data:    []byte(id[:]),
// 	}

// 	filter := bson.D{{
// 		Key:   "_id",
// 		Value: binID,
// 	}}

// 	findChan := make(chan error)
// 	delChan := make(chan error)

// 	var f file.File
// 	go func() {
// 		findChan <- db.FindOne(
// 			ctx, bson.D{{
// 				Key:   "_id",
// 				Value: binID,
// 			}},
// 		).Decode(&f)
// 	}()

// 	go func() {
// 		_, err := db.DeleteOne(ctx, filter)
// 		delChan <- err
// 	}()

// 	select {
// 	case err := <-delChan:
// 		if err != nil {
// 			<-findChan
// 			log.Println(err)
// 			return file.ErrInternal
// 		}
// 	case err := <-findChan:
// 		if err != nil {
// 			<-delChan
// 			log.Println(err)

// 			if err == mongo.ErrNoDocuments {
// 				return file.ErrFileNotFoundByID
// 			}

// 			return file.ErrInternal
// 		}
// 	}

// 	select {
// 	case err := <-delChan:
// 		if err != nil {
// 			<-findChan
// 			log.Println(err)
// 			return file.ErrInternal
// 		}
// 	case err := <-findChan:
// 		if err != nil {
// 			<-delChan
// 			log.Println(err)

// 			if err == mongo.ErrNoDocuments {
// 				return file.ErrFileNotFoundByID
// 			}

// 			return file.ErrInternal
// 		}
// 	}

// 	return nil
// }
