package file

import (
	"context"

	"github.com/google/uuid"
)

type FindByIDRequest struct {
	ID uuid.UUID
}

type FindByIDResponse struct {
	Metadata *File
}

func (s *Service) FindByID(ctx context.Context, req FindByIDRequest) (FindByIDResponse, error) {
	f, err := s.Repo.FindByID(ctx, req.ID)
	if err != nil {
		return FindByIDResponse{}, err
	}

	return FindByIDResponse{f}, nil
}
