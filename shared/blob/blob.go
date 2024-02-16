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
	ErrInternal         = errors.New("blob: error while communicating with storage")
	ErrInvalidEnv       = errors.New(`blob: application environment mode ("APP_ENV") is missing or not set`)
	ErrReceivedSignal   = errors.New("blob: signal received to abort the operation")
	ErrEmptyContentType = errors.New("blob: Content-Type is empty")
	ErrBucketNotFound   = errors.New("blob: bucket was not found")
	ErrEmptyBucket      = errors.New("blob: bucket name is empty")
	ErrBlobNotFound     = errors.New("blob: key was not found")
	ErrEmptyBlobKey     = errors.New("blob: key is empty")
	ErrNilReader        = errors.New("blob: io.Reader is nil")
	ErrNilLocation      = errors.New(`blob: location ("fmt.Stringer") is nil`)
)
