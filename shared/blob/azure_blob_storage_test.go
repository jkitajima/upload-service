package blob

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"gocloud.dev/blob"
)

func TestNewAzureBlobStorage(t *testing.T) {
	if os.Getenv("AZURE_STORAGE_ACCOUNT") == "" || os.Getenv("AZURE_STORAGE_KEY") == "" {
		t.Skip("blob: azure_blob_storage: test_new_azure_blob_storage: required environment vars (`AZURE_STORAGE_ACCOUNT`, `AZURE_STORAGE_KEY`) not set. could not connect to a test blob storage")
	}

	cases := map[string]struct {
		inSetenv  func(key, value string)
		inAccount string
		inKey     string
		outStg    Storager
		outErr    error
	}{
		"new repo": {t.Setenv, "devstoreaccount1", "Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==", &azure{domain: "https://devstoreaccount1.blob.core.windows.net/"}, nil},
		"missing environment variable `AZURE_STORAGE_ACCOUNT`": {t.Setenv, "", "Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==", nil, ErrAccountEnvVar},
		"missing environment variable `AZURE_STORAGE_KEY`":     {t.Setenv, "devstoreaccount1", "", nil, ErrKeyEnvVar},
	}

	for key, testcase := range cases {
		t.Run(key, func(t *testing.T) {
			testcase.inSetenv("AZURE_STORAGE_ACCOUNT", testcase.inAccount)
			testcase.inSetenv("AZURE_STORAGE_KEY", testcase.inKey)

			blobstg, err := NewAzureBlobStorage()
			if err != nil {
				if err != testcase.outErr {
					t.Errorf("blob: azure_blob_storage: test_new_azure_blob_storage: error interface mismatch (result = %q, expected = %q)\n", err, testcase.outErr)
				}
				return
			}

			switch {
			case blobstg.String() != testcase.outStg.String():
				t.Errorf("blob: azure_blob_storage: test_new_azure_blob_storage: `Storager` interface mismatch (result = %q, expected = %q)\n", blobstg, testcase.outStg)
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
