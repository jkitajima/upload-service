package mongo

import (
	"context"
	"testing"
	"time"

	"upload/pkg/file"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestFindByID(t *testing.T) {
	const timeout = 15 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	t.Cleanup(func() { cancel() })

	uri := "mongodb://localhost:27017"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		t.Skipf("file: repo: mongo: failed to connect with test db: %v", err)
	}
	t.Cleanup(func() { client.Disconnect(ctx) })

	db := client.Database("Upload")
	coll := db.Collection("Files-find_by_id")
	coll.Drop(ctx)
	t.Cleanup(func() { coll.Drop(ctx) })
	fileColl := NewRepo(coll)

	id, _ := uuid.Parse("6bd1edba-bfe6-4a18-b2f3-db1dc4043b32")
	f := file.File{ID: id, UploaderID: "0", CompanyID: "1", Name: "photo.png", Extension: ".png", ContentType: "image/png", Size: 256, StorageLocation: "https://storage.blob.com/" + id.String(), Description: "some random photo", SubmittedAt: time.Now()}
	fileColl.Insert(ctx, &f)

	cases := map[string]struct {
		inID    uuid.UUID
		outFile *file.File
		outErr  error
	}{
		"fetching a file should increment the requested counter": {id, &f, nil},
		"file not found": {uuid.New(), &f, file.ErrNotFoundByID},
	}

	for key, testcase := range cases {
		t.Run(key, func(t *testing.T) {
			f, err := fileColl.FindByID(ctx, testcase.inID)
			if err != nil {
				if err != testcase.outErr {
					t.Errorf("file: repo: mongo: test_find_by_id: error mismatch (err = %v, expected = %v)\n", err, testcase.outErr)
				}
				return
			}

			switch {
			case err != testcase.outErr:
				t.Errorf("file: repo: mongo: test_find_by_id: error mismatch (err = %v, expected = %v)\n", err, testcase.outErr)
			case f.ID != testcase.inID:
				t.Errorf("file: repo: mongo: test_find_by_id: ID mismatch (result = %v, expected = %v)\n", err, testcase.outErr)
			case f.TimesRequested != 1:
				t.Errorf("file: repo: mongo: test_find_by_id: times request counter mismatch (result = %v, expected = %v)\n", err, testcase.outErr)
			}
		})
	}
}
