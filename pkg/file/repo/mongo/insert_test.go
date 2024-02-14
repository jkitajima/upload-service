package mongo

import (
	"context"
	"os"
	"testing"
	"time"

	"upload/pkg/file"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestInsert(t *testing.T) {
	const (
		location   = "upload/pkg/file/repo/mongo/insert_test.go"
		function   = "TestInsert"
		timeout    = 15 * time.Second
		mongodbURI = "MONGODB_URI"
		dbName     = "upload-test"
		collName   = "files-insert"
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	t.Cleanup(func() { cancel() })

	uri := os.Getenv(mongodbURI)
	if uri == "" {
		t.Skipf(`
	location: %s
	func: %s
	msg: environment variable %q is either empty or missing`,
			location, function, mongodbURI)
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		t.Skipf(`
	location: %s
	func: %s
	msg: failed to connect with test db
	errmsg: %v`,
			location, function, err)
	}
	t.Cleanup(func() { client.Disconnect(ctx) })

	db := client.Database(dbName)
	coll := db.Collection(collName)
	coll.Drop(ctx)
	t.Cleanup(func() { coll.Drop(ctx) })
	fileColl := NewRepo(coll)

	id, _ := uuid.Parse("6bd1edba-bfe6-4a18-b2f3-db1dc4043b32")
	f := file.File{ID: id, UploaderID: "0", CompanyID: "1", Name: "photo.png", Extension: ".png", ContentType: "image/png", Size: 256, StorageLocation: "https://storage.blob.com/" + id.String(), Description: "some random photo", SubmittedAt: time.Now()}

	cases := map[string]struct {
		in  *file.File
		out error
	}{
		"basic insert":   {&f, nil},
		"duplicated key": {&f, file.ErrInsertDuplicatedKey},
	}

	for key, testcase := range cases {
		t.Run(key, func(t *testing.T) {
			if err := fileColl.Insert(ctx, testcase.in); err != testcase.out {
				t.Errorf("file: repo: mongo: test_insert: %v\n", err)
			}
		})
	}
}
