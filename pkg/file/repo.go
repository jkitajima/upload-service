package file

import (
	"context"

	"github.com/google/uuid"
)

type repo interface {
	Insert(context.Context, *File) error
	FindByID(context.Context, uuid.UUID) (*File, error)
	UpdateByID(context.Context, uuid.UUID, *File) error
	DeleteByID(context.Context, uuid.UUID) error
}
