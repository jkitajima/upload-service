package file

import (
	"context"

	"github.com/google/uuid"
)

type repo interface {
	Insert(context.Context, *File) error
	FindByID(context.Context, uuid.UUID) (*File, error)
	Update(context.Context, uuid.UUID, *File) error
	Delete(context.Context, uuid.UUID) error
}
