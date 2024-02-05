package blob

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/azureblob"
)

type azure struct {
	scheme string
	domain string
}

func NewAzureBlobStorage() (Storager, error) {
	domain := os.Getenv("AZURE_STORAGE_ACCOUNT")
	key := os.Getenv("AZURE_STORAGE_KEY")

	switch {
	case domain == "":
		return nil, ErrAccountEnvVar
	case key == "":
		return nil, ErrKeyEnvVar
	}

	return &azure{
		scheme: "azblob://",
		domain: fmt.Sprintf("https://%s.blob.core.windows.net/", domain),
	}, nil
}

func (az *azure) String() string { return az.domain }

func (az *azure) Upload(ctx context.Context, bucket, key string, r io.Reader, opts *blob.WriterOptions) error {
	switch {
	case r == nil:
		err := ErrNilReader
		log.Printf("azure blob storage: upload: blob writer options: %v\n", err)
		return err
	case opts.ContentType == "":
		err := ErrEmptyContentType
		log.Printf("azure blob storage: upload: blob writer options: %v\n", err)
		return err
	case bucket == "":
		err := ErrEmptyBucket
		log.Printf("azure blob storage: upload: blob writer options: %v\n", err)
		return err
	case key == "":
		err := ErrEmptyBlobKey
		log.Printf("azure blob storage: upload: blob writer options: %v\n", err)
		return err
	}

	buck, err := blob.OpenBucket(ctx, az.scheme+bucket)
	if err != nil {
		log.Printf("azure blob storage: upload: bucket opening: %v\n", err)
		return ErrInternal
	}
	defer buck.Close()

	if err := buck.Upload(ctx, key, r, opts); err != nil {
		log.Printf("blob: azure blob storage: upload: bucket upload: %v\n", err)
		switch {
		case bloberror.HasCode(err, bloberror.ContainerNotFound):
			err = ErrBucketNotFound
		default:
			err = ErrInternal
		}
		return err
	}

	return nil
}

func (az *azure) Delete(ctx context.Context, bucket, item string) error {
	switch {
	case bucket == "":
		err := ErrEmptyBucket
		log.Printf("blob: azure blob storage: delete: %v\n", err)
		return err
	case item == "":
		err := ErrEmptyBlobKey
		log.Printf("blob: azure blob storage: delete: %v\n", err)
		return err
	}

	buck, err := blob.OpenBucket(ctx, az.scheme+bucket)
	if err != nil {
		log.Printf("azure blob storage: delete: bucket opening: %v\n", err)
		return ErrInternal
	}
	defer buck.Close()

	if err := buck.Delete(ctx, item); err != nil {
		log.Printf("blob: azure blob storage: delete: bucket delete: %v\n", err)
		switch {
		case bloberror.HasCode(err, bloberror.ContainerNotFound):
			err = ErrBucketNotFound
		case bloberror.HasCode(err, bloberror.BlobNotFound):
			err = ErrBlobNotFound
		default:
			err = ErrInternal
		}
		return err
	}

	return nil
}

func (az *azure) KillZombie(location fmt.Stringer) error {
	if location == nil {
		return ErrNilLocation
	}

	bucket := location.(*Location).Bucket
	if bucket == "" {
		return ErrEmptyBucket
	}

	key := location.(*Location).Key
	if key == "" {
		return ErrEmptyBlobKey
	}

	buck, err := blob.OpenBucket(context.TODO(), az.scheme+bucket)
	if err != nil {
		log.Println(err)
		return ErrInternal
	}
	defer buck.Close()

	if err := buck.Delete(context.TODO(), key); err != nil {
		log.Printf("blob: azure blob storage: delete: bucket delete: %v\n", err)
		switch {
		case bloberror.HasCode(err, bloberror.ContainerNotFound):
			err = ErrBucketNotFound
		case bloberror.HasCode(err, bloberror.BlobNotFound):
			err = ErrBlobNotFound
		default:
			err = ErrInternal
		}
		return err
	}

	return nil
}
