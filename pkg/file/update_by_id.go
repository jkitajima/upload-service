package file

import (
	"context"

	"github.com/google/uuid"
)

func Update(ctx context.Context, r repo, id uuid.UUID, f *File) error {
	if err := r.Update(ctx, id, f); err != nil {
		return err
	}

	return nil
}
