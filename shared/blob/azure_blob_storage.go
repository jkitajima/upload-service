package blob

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"upload/shared/zombiekiller"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/azureblob"
	"gocloud.dev/gcerrors"
)

type azure struct {
	scheme string
	domain string
}

func NewAzureBlobStorage() (Storager, error) {
	domain := os.Getenv("AZURE_STORAGE_ACCOUNT")
	if domain == "" {
		return nil, errors.New("environment variable `AZURE_STORAGE_ACCOUNT` is either empty or not set")
	}

	return &azure{
		scheme: "azblob://",
		domain: fmt.Sprintf("https://%s.blob.core.windows.net/", domain),
	}, nil
}

func (az *azure) String() string { return az.domain }

func (az *azure) Upload(ctx context.Context, bucket, key string, r io.Reader, opts *blob.WriterOptions) error {
	if opts.ContentType == "" {
		err := ErrEmptyContentType
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
		log.Printf("azure blob storage: upload: bucket upload: %v\n", err)
		return ErrInternal
	}

	return nil
}

func (az *azure) Delete(ctx context.Context, bucket, item string) error {
	buck, err := blob.OpenBucket(ctx, az.scheme+bucket)
	if err != nil {
		log.Printf("azure blob storage: delete: bucket opening: %v\n", err)
		return ErrInternal
	}
	defer buck.Close()

	if err := buck.Delete(ctx, item); err != nil {
		log.Printf("azure blob storage: upload: bucket delete: %v\n", err)

		errcode := gcerrors.Code(err)
		if errcode == gcerrors.NotFound {
			return ErrNotFound
		}

		return ErrInternal
	}

	return nil
}

func (az *azure) KillZombie(location fmt.Stringer) error {
	bucket := location.(*Location).Bucket
	key := location.(*Location).Key

	buck, err := blob.OpenBucket(context.TODO(), az.scheme+bucket)
	if err != nil {
		log.Println(err)
		return ErrInternal
	}
	defer buck.Close()

	if err := buck.Delete(context.TODO(), key); err != nil {
		log.Println(err)

		errcode := gcerrors.Code(err)
		if errcode == gcerrors.NotFound {
			return zombiekiller.ErrNotFound
		}

		return zombiekiller.ErrInternal
	}

	return nil
}
