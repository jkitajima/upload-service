package httphandler

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"upload/pkg/file"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestHandleFileFindByID(t *testing.T) {
	const timeout = 15 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	t.Cleanup(func() { cancel() })

	uri := "mongodb://localhost:27017"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		t.Skipf("file: httphandler: test_handle_file_find_by_id: failed to connect with test db: %v\n", err)
	}
	t.Cleanup(func() { client.Disconnect(ctx) })

	db := client.Database("Upload")
	coll := db.Collection("Files-handle_file_find_by_id")
	coll.Drop(ctx)
	t.Cleanup(func() { coll.Drop(ctx) })

	srv := NewServer(coll, nil, nil)
	id, _ := uuid.Parse("6bd1edba-bfe6-4a18-b2f3-db1dc4043b32")
	f := file.File{ID: id, UploaderID: "0", CompanyID: "1", Name: "photo.png", Extension: ".png", ContentType: "image/png", Size: 256, StorageLocation: "https://storage.blob.com/" + id.String(), Description: "some random photo", SubmittedAt: time.Now()}
	srv.db.Insert(ctx, &f)
	get := http.MethodGet

	cases := map[string]struct {
		reqMethod string
		reqTarget string
		reqBody   io.Reader
		respCode  int
	}{
		"fetch an existent file by id should be ok":            {get, "/" + id.String(), nil, http.StatusOK},
		"fetching a non existent file should return not found": {get, "/0bd1edba-3333-4444-7777-db1dc4043b39", nil, http.StatusNotFound},
		"fetching a resource with an invalid uuid":             {get, "/1996", nil, http.StatusBadRequest},
	}

	for key, testcase := range cases {
		t.Run(key, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(testcase.reqMethod, testcase.reqTarget, testcase.reqBody)
			srv.ServeHTTP(w, r)
			require.Equal(t, testcase.respCode, w.Result().StatusCode)
		})
	}
}
