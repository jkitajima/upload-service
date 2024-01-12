package blob

import (
	"context"
	"log"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/azureblob"
	"gocloud.dev/gcerrors"
)

const scheme = "azblob://"

func Upload(ctx context.Context, bucket, item string, data []byte, opts *blob.WriterOptions) error {
	buck, err := blob.OpenBucket(ctx, scheme+bucket)
	if err != nil {
		log.Println(err)
		return ErrInternal
	}
	defer buck.Close()

	w, err := buck.NewWriter(ctx, item, opts)
	if err != nil {
		log.Println(err)
		return ErrInternal
	}

	_, writeErr := w.Write(data)
	closeErr := w.Close()

	if writeErr != nil {
		log.Println(err)
		return ErrInternal
	}
	if closeErr != nil {
		log.Println(err)
		return ErrInternal
	}

	return nil
}

func Delete(ctx context.Context, bucket, item string) error {
	buck, err := blob.OpenBucket(ctx, scheme+bucket)
	if err != nil {
		log.Println(err)
		return ErrInternal
	}
	defer buck.Close()

	if err := buck.Delete(ctx, item); err != nil {
		log.Println(err)

		errcode := gcerrors.Code(err)
		if errcode == gcerrors.NotFound {
			return ErrNotFound
		}

		return ErrInternal
	}

	return nil
}
