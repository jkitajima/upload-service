package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	fileServer "upload/pkg/file/httphandler"
	"upload/shared/blob"
	"upload/shared/composer"
	"upload/shared/zombiekiller"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	env := flag.String("env", "unset", "set which environment to load variables")
	flag.Parse()

	if err := loadenv(*env); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func loadenv(env string) error {
	switch env {
	case "unset":
		fmt.Printf("Program executed without setting an environment. Using default option: %q.\n", "dev")
		fallthrough
	case "dev":
		env = "dev.env"
	case "test":
		env = "test.env"
	default:
		return fmt.Errorf("invalid environment. valid options are: [%q, %q]", "dev", "test")
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	env = filepath.Join(wd, env)
	if err := godotenv.Overload(env); err != nil {
		return err
	}

	return nil
}

func run() error {
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
	port := os.Getenv("APP_PORT")
	if port == "" {
		return errors.New("environment variable `APP_PORT` is null or non-existent")
	}

	// init zombie killer
	doneChan := make(chan any)
	defer close(doneChan)
	const thrashBuffer = 1 << 10 * 1 // 1024 * (servers count)
	thrashChan := make(chan zombiekiller.KillOperation, thrashBuffer)
	go zombiekiller.ListenForKillOperations(doneChan, thrashChan)

	// init servers
	blobstg, err := blob.NewAzureBlobStorage()
	if err != nil {
		return err
	}

	srv := composer.NewComposer()
	file := fileServer.NewServer(dbClient.Collection("files"), blobstg, thrashChan)
	if err := srv.Compose(file); err != nil {
		return err
	}

	log.Printf("server listening on port %s...\n", port)
	log.Fatalln(http.ListenAndServe(":"+port, srv))

	return nil
}
