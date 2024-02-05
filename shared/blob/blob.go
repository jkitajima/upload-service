package blob

import (
	"context"
	"errors"
	"fmt"
	"io"

	"upload/shared/zombiekiller"

	"gocloud.dev/blob"
)

type Storager interface {
	fmt.Stringer
	Upload(ctx context.Context, bucket, key string, r io.Reader, opts *blob.WriterOptions) error
	Delete(ctx context.Context, bucket, key string) error
	zombiekiller.ZombieKiller
}

type Location struct{ Bucket, Key string }

func (l *Location) String() string {
	return fmt.Sprintf("%s/%s", l.Bucket, l.Key)
}

var (
	ErrInternal         = errors.New("error while communicating with blob storage")
	ErrReceivedSignal   = errors.New("blob received a signal to abort the operation")
	ErrEmptyContentType = errors.New("blob Content-Type is empty")
	ErrBucketNotFound   = errors.New("bucket was not found")
	ErrEmptyBucket      = errors.New("bucket name is empty")
	ErrBlobNotFound     = errors.New("blob was not found")
	ErrEmptyBlobKey     = errors.New("blob key is empty")
	ErrNilReader        = errors.New("io.Reader is nil")
	ErrAccountEnvVar    = errors.New("environment variable `AZURE_STORAGE_ACCOUNT` is either empty or not set")
	ErrKeyEnvVar        = errors.New("environment variable `AZURE_STORAGE_KEY` is either empty or not set")
)
