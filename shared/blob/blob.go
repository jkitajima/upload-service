package blob

import (
	"context"
	"errors"
	"io"

	"gocloud.dev/blob"
)

type Storager interface {
	Upload(ctx context.Context, bucket, key string, r io.Reader, opts *blob.WriterOptions) error
	Delete(ctx context.Context, bucket, key string) error
}

var (
	ErrInternal         = errors.New("error while communicating with blob storage")
	ErrReceivedSignal   = errors.New("blob received a signal to abort the operation")
	ErrEmptyContentType = errors.New("blob Content-Type is empty")
	ErrNotFound         = errors.New("blob was not found")
)
