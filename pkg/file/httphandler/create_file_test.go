package httphandler

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"upload/shared/blob"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestHandleFileCreate(t *testing.T) {
	const timeout = 15 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	t.Cleanup(func() { cancel() })

	uri := "mongodb://localhost:27017"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		t.Skipf("file: httphandler: test_handle_file_create: failed to connect with test db: %v\n", err)
	}
	t.Cleanup(func() { client.Disconnect(ctx) })

	db := client.Database("upload")
	coll := db.Collection("files.handle_file_create")
	coll.Drop(ctx)
	t.Cleanup(func() { coll.Drop(ctx) })

	if os.Getenv("AZURE_STORAGE_ACCOUNT") == "" || os.Getenv("AZURE_STORAGE_KEY") == "" {
		t.Skip("file: httphandler: test_handle_file_create: required environment vars (`AZURE_STORAGE_ACCOUNT`, `AZURE_STORAGE_KEY`) not set. could not connect to a test blob storage")
	}

	blobstg, err := blob.NewAzureBlobStorage()
	if err != nil {
		t.Skipf("file: httphandler: test_handle_file_create: failed to create a blob storage for testing (err = %q)\n", err)
	}

	s := NewServer(nil, coll, blobstg, nil)
	ts := httptest.NewServer(s)
	t.Cleanup(func() { ts.Close() })
	baseURL := ts.URL
	httpclient := ts.Client()

	var b bytes.Buffer
	wr := multipart.NewWriter(&b)
	fwr, _ := wr.CreateFormFile("file", "data.txt")
	fwr.Write([]byte("test data from data.txt file\n"))
	wr.WriteField("uploaderId", "123")
	wr.WriteField("companyId", "444")
	wr.WriteField("description", "some random text")
	wr.Close()
	ct := wr.FormDataContentType()

	cases := map[string]struct {
		reqTarget      string
		reqBody        io.Reader
		reqContentType string
		respCode       int
	}{
		"basic insert": {baseURL + "/", &b, ct, http.StatusCreated},
	}

	for key, testcase := range cases {
		t.Run(key, func(t *testing.T) {
			resp, err := httpclient.Post(testcase.reqTarget, testcase.reqContentType, testcase.reqBody)
			if err != nil {
				t.Errorf("file: httphandler: test_handle_file_create: post failed (err = %v)\n", err)
			}

			require.Equal(t, testcase.respCode, resp.StatusCode)
		})
	}
}
