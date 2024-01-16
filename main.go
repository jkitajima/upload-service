package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	fileServer "upload/pkg/file/httphandler"
	"upload/shared/composer"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := exec(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func exec() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	// connect to mongodb
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return errors.New("environment variable `MONGODB_URI` is null or non-existent")
	}

	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	defer func() {
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	dbName := os.Getenv("MONGODB_NAME")
	if dbName == "" {
		return errors.New("environment variable `MONGODB_NAME` is null or non-existent")
	}
	dbClient := mongoClient.Database(dbName)

	// compose pkg servers
	port := os.Getenv("PORT_NUMBER")
	if port == "" {
		return errors.New("environment variable `PORT_NUMBER` is null or non-existent")
	}

	srv := composer.NewComposer()
	file := fileServer.NewServer(dbClient.Collection("files"), "azblob://")
	if err := srv.Compose(file); err != nil {
		return err
	}

	log.Printf("server listening on port %s...\n", port)
	log.Fatalln(http.ListenAndServe(":"+port, srv))

	return nil
}
