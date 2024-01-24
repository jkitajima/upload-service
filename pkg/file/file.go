package file

import (
	"context"
	"errors"
	"time"

	"upload/shared/blob"
	"upload/shared/zombiekiller"

	"github.com/google/uuid"
)

type File struct {
	ID              uuid.UUID
	UploaderID      string
	CompanyID       string
	Name            string
	Extension       string
	ContentType     string
	Size            uint
	StorageLocation string
	TimesRequested  uint
	Description     string
	SubmittedAt     time.Time
	UpdatedAt       time.Time
	UploadedAt      time.Time
}

type Service struct {
	Repo   Repoer
	Blob   blob.Storager
	Thrash chan<- zombiekiller.KillOperation
}

type Repoer interface {
	Insert(context.Context, *File) error
	FindByID(context.Context, uuid.UUID) (*File, error)
	UpdateByID(context.Context, uuid.UUID, *File) error
	DeleteByID(context.Context, uuid.UUID) error
	zombiekiller.ZombieKiller
}

var (
	ErrInternal            = errors.New("the file service encountered an unexpected condition that prevented it from fulfilling the request")
	ErrRepoReceivedSignal  = errors.New("repository received a signal to abort the operation")
	ErrNotFoundByID        = errors.New("could not find any file with provided ID")
	ErrEmptyContentType    = errors.New("file Content-Type is empty")
	ErrEmptyBucketName     = errors.New("provided bucket name is empty")
	ErrInvalidUUID         = errors.New("provided file UUID is invalid")
	ErrInsertDuplicatedKey = errors.New("insertion failed because a file with the provided key already exists")
	ErrTimeout             = errors.New("requested service operation timed out")
)
