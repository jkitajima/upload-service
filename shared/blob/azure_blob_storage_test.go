package blob

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"gocloud.dev/blob"
)

func TestNewAzureBlobStorage(t *testing.T) {
	cases := map[string]struct {
		err error
	}{
		"new repo": {nil},
	}

	for key, testcase := range cases {
		t.Run(key, func(t *testing.T) {
			_, err := NewAzureBlobStorage()
			if err != nil {
				if err != testcase.err {
					t.Errorf("blob: azure_blob_storage: test_new_azure_blob_storage: error interface mismatch (result = %q, expected = %q)\n", err, testcase.err)
				}
				return
			}
		})
	}
}

func TestUpload(t *testing.T) {
	const timeout = 15 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	t.Cleanup(func() { cancel() })

	blobstg, err := NewAzureBlobStorage()
	if err != nil {
		t.Skipf("blob: azure_blob_storage: test_new_azure_blob_storage: failed to create a blob storage for testing (err = %q)\n", err)
	}

	cases := map[string]struct {
		inBucket, inKey string
		inReader        io.Reader
		inOpts          *blob.WriterOptions
		outErr          error
	}{
		"basic upload":        {"company", "file1.txt", strings.NewReader("file1 test data"), &blob.WriterOptions{ContentType: "text/plain"}, nil},
		"non-existent bucket": {"c0mpany", "file2.txt", strings.NewReader("file2 test data"), &blob.WriterOptions{ContentType: "text/plain"}, ErrBucketNotFound},
		"missing key":         {"company", "", strings.NewReader("file3 test data"), &blob.WriterOptions{ContentType: "text/plain"}, ErrEmptyBlobKey},
		"missing bucket name": {"", "file4.txt", strings.NewReader("file4 test data"), &blob.WriterOptions{ContentType: "text/plain"}, ErrEmptyBucket},
		"empty content-type":  {"company", "file5.txt", strings.NewReader("file5 test data"), &blob.WriterOptions{ContentType: ""}, ErrEmptyContentType},
		"empty/nil reader":    {"company", "file6.txt", nil, &blob.WriterOptions{ContentType: "text/plain"}, ErrNilReader},
	}

	for key, testcase := range cases {
		t.Run(key, func(t *testing.T) {
			if err := blobstg.Upload(ctx, testcase.inBucket, testcase.inKey, testcase.inReader, testcase.inOpts); err != testcase.outErr {
				t.Errorf("blob: azure_blob_storage: test_new_azure_blob_storage: error mismatched (err = %q, expected = %q)\n", err, testcase.outErr)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	const timeout = 15 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	t.Cleanup(func() { cancel() })

	blobstg, err := NewAzureBlobStorage()
	if err != nil {
		t.Skipf("blob: azure_blob_storage: test_new_azure_blob_storage: failed to create a blob storage for testing (err = %q)\n", err)
	}

	cases := map[string]struct {
		inBucket, inKey string
		outErr          error
	}{
		"empty bucket name":       {"", "file1.txt", ErrEmptyBucket},
		"empty blob key":          {"company", "", ErrEmptyBlobKey},
		"bucket does not exist":   {"c0mpany", "file1.txt", ErrBucketNotFound},
		"blob key does not exist": {"company", "file2.txt", ErrBlobNotFound},
		"basic delete":            {"company", "file1.txt", nil},
	}

	for key, testcase := range cases {
		t.Run(key, func(t *testing.T) {
			if err := blobstg.Delete(ctx, testcase.inBucket, testcase.inKey); err != testcase.outErr {
				t.Errorf("blob: azure_blob_storage: test_new_azure_blob_storage: delete: error mismatched (err = %q, expected = %q)\n", err, testcase.outErr)
			}
		})
	}
}

func TestKillZombie(t *testing.T) {
	blobstg, err := NewAzureBlobStorage()
	if err != nil {
		t.Skipf("blob: azure_blob_storage: test_new_azure_blob_storage: failed to create a blob storage for testing (err = %q)\n", err)
	}

	blobstg.Upload(context.TODO(), "company", "file1.txt", strings.NewReader("test data"), &blob.WriterOptions{ContentType: "text/plain"})

	cases := map[string]struct {
		inLocation fmt.Stringer
		outErr     error
	}{
		"nil stringer interface":  {nil, ErrNilLocation},
		"empty bucket name":       {&Location{"", "file1.txt"}, ErrEmptyBucket},
		"empty blob key":          {&Location{"company", ""}, ErrEmptyBlobKey},
		"bucket does not exist":   {&Location{"doesnotexist", "file1.txt"}, ErrBucketNotFound},
		"blob key does not exist": {&Location{"company", "doesnotexist.txt"}, ErrBlobNotFound},
		"basic kill":              {&Location{"company", "file1.txt"}, nil},
	}

	for key, testcase := range cases {
		t.Run(key, func(t *testing.T) {
			if err := blobstg.KillZombie(testcase.inLocation); err != testcase.outErr {
				t.Errorf("blob: azure_blob_storage: test_new_azure_blob_storage: kill zombie: error mismatched (err = %q, expected = %q)\n", err, testcase.outErr)
			}
		})
	}
}
