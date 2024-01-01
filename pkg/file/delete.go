package file

import (
	"context"

	"github.com/google/uuid"
)

func Delete(ctx context.Context, r repo, id uuid.UUID) error {
	if err := r.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
