package blob

import (
	"context"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/azureblob"
)

const scheme = "azblob://"

func Upload(ctx context.Context, bucket, item string, data []byte, opts *blob.WriterOptions) error {
	buck, err := blob.OpenBucket(ctx, scheme+bucket)
	if err != nil {
		return err
	}
	defer buck.Close()

	w, err := buck.NewWriter(ctx, item, opts)
	if err != nil {
		return err
	}

	_, writeErr := w.Write(data)
	closeErr := w.Close()

	if writeErr != nil {
		return err
	}
	if closeErr != nil {
		return err
	}

	return nil
}

func Delete(ctx context.Context, bucket, item string) error {
	buck, err := blob.OpenBucket(ctx, scheme+bucket)
	if err != nil {
		return err
	}
	defer buck.Close()

	if err := buck.Delete(ctx, item); err != nil {
		return err
	}

	return nil
}
