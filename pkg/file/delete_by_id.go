package file

import (
	"context"
	"log"
	"sync"

	"upload/shared/blob"
	"upload/shared/zombiekiller"

	"github.com/google/uuid"
)

type DeleteByIDRequest struct {
	ID     uuid.UUID
	Bucket string
}

func (s *Service) DeleteByID(ctx context.Context, req DeleteByIDRequest) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	blobChan := make(chan error)
	go func() {
		blobChan <- s.Blob.Delete(ctx, req.Bucket, req.ID.String())
	}()

	repoChan := make(chan error)
	go func() {
		repoChan <- s.Repo.DeleteByID(ctx, req.ID)
	}()

	sentinel := struct {
		mu          sync.Mutex
		responded   bool
		errored     bool
		errResponse error
	}{}

	for i := 0; i < 2; i++ {
		select {
		case err := <-repoChan:
			sentinel.mu.Lock()
			if err != nil {
				cancel()
				log.Printf("file: service: delete_by_id: repo: %v\n", err)
				if sentinel.responded && !sentinel.errored {
					s.Thrash <- zombiekiller.KillOperation{Killer: s.Repo, Target: req.ID}
				}
				sentinel.errored = true
				sentinel.errResponse = err
			}
			sentinel.responded = true
			sentinel.mu.Unlock()
		case err := <-blobChan:
			sentinel.mu.Lock()
			if err != nil {
				cancel()
				log.Printf("file: service: delete_by_id: blob: %v\n", err)
				if sentinel.responded && !sentinel.errored {
					s.Thrash <- zombiekiller.KillOperation{Killer: s.Blob, Target: &blob.Location{Bucket: req.Bucket, Key: req.ID.String()}}
				}
				sentinel.errored = true
				sentinel.errResponse = err
			}
			sentinel.responded = true
			sentinel.mu.Unlock()
		}
	}

	if sentinel.errored {
		return sentinel.errResponse
	}
	return nil
}
