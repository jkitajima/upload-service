package file

import (
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gocloud.dev/blob"
)

type repoerMock struct{ mock.Mock }

func (r *repoerMock) Insert(ctx context.Context, f *File) error {
	args := r.Called(f)
	return args.Error(0)
}

func (r *repoerMock) FindByID(ctx context.Context, id uuid.UUID) (*File, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(*File), args.Error(1)
}

func (r *repoerMock) UpdateByID(ctx context.Context, id uuid.UUID, f *File) error {
	args := r.Called(ctx, id, f)
	return args.Error(0)
}

func (r *repoerMock) DeleteByID(ctx context.Context, id uuid.UUID) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}

func (r *repoerMock) KillZombie(key fmt.Stringer) error {
	args := r.Called(key)
	return args.Error(0)
}

type storagerMock struct{ mock.Mock }

func (s *storagerMock) String() string {
	return "storager mock"
}

func (s *storagerMock) Upload(ctx context.Context, bucket, key string, r io.Reader, opts *blob.WriterOptions) error {
	args := s.Called(bucket, key, r, opts)
	return args.Error(0)
}

func (s *storagerMock) Delete(ctx context.Context, bucket, key string) error {
	args := s.Called(ctx, bucket, key)
	return args.Error(0)
}

func (r *storagerMock) KillZombie(key fmt.Stringer) error {
	args := r.Called(key)
	return args.Error(0)
}
