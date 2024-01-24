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

func TestUpdateByID(t *testing.T) {
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
	coll := db.Collection("Files-update_by_id")
	coll.Drop(ctx)
	t.Cleanup(func() { coll.Drop(ctx) })
	fileColl := NewRepo(coll)

	id, _ := uuid.Parse("6bd1edba-bfe6-4a18-b2f3-db1dc4043b32")
	f := file.File{ID: id, UploaderID: "0", CompanyID: "1", Name: "photo.png", Extension: ".png", ContentType: "image/png", Size: 256, StorageLocation: "https://storage.blob.com/" + id.String(), Description: "some random photo", SubmittedAt: time.Now()}
	fileColl.Insert(ctx, &f)

	cases := map[string]struct {
		inID   uuid.UUID
		inFile *file.File
		outErr error
	}{
		"basic update":                  {id, &file.File{UploaderID: "444", CompanyID: "555", Description: "Nova descrição aqui"}, nil},
		"requested file does not exist": {uuid.New(), &file.File{UploaderID: "444", CompanyID: "555", Description: "Nova descrição aqui"}, file.ErrNotFoundByID},
	}

	for key, testcase := range cases {
		t.Run(key, func(t *testing.T) {
			if err := fileColl.UpdateByID(ctx, testcase.inID, testcase.inFile); err != testcase.outErr {
				t.Errorf("file: repo: mongo: test_update_by_id: error mismatch (result = %q, expected = %q)\n", err, testcase.outErr)
				return
			}

			f, err := fileColl.FindByID(ctx, testcase.inID)
			if err != nil {
				if err != testcase.outErr {
					t.Errorf("file: repo: mongo: test_update_by_id: error mismatch (err = %v, expected = %v)\n", err, testcase.outErr)
				}
				return
			}

			switch {
			case f.ID != testcase.inID:
				t.Errorf("file: repo: mongo: test_update_by_id: updated file was requested and mismatched (err = %v)\n", nil)
				return
			case f.UploaderID != testcase.inFile.UploaderID:
				t.Errorf("file: repo: mongo: test_update_by_id: uploaderID mismatch (result = %q, expected = %q)\n", f.UploaderID, testcase.inFile.UploaderID)
			case f.CompanyID != testcase.inFile.CompanyID:
				t.Errorf("file: repo: mongo: test_update_by_id: companyID mismatch (result = %q, expected = %q)\n", f.CompanyID, testcase.inFile.CompanyID)
			case f.Description != testcase.inFile.Description:
				t.Errorf("file: repo: mongo: test_update_by_id: description mismatch (result = %q, expected = %q)\n", f.Description, testcase.inFile.Description)
			}
		})
	}
}
