package file

import (
	"context"

	"github.com/google/uuid"
)

func FindByID(ctx context.Context, r repo, id uuid.UUID) (*File, error) {
	f, err := r.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return f, nil
}
