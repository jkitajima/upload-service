package file

import (
	"context"

	"github.com/google/uuid"
)

type UpdateByIDRequest struct {
	ID       uuid.UUID
	Metadata *File
}

type UpdateByIDResponse struct {
	Metadata *File
}

func (s *Service) UpdateByID(ctx context.Context, req UpdateByIDRequest) (UpdateByIDResponse, error) {
	err := s.Repo.UpdateByID(ctx, req.ID, req.Metadata)
	if err != nil {
		return UpdateByIDResponse{}, err
	}

	return UpdateByIDResponse{Metadata: req.Metadata}, nil
}
