package file

import (
	"context"
	"io"
	"log"
	"sync"
	"upload/shared/blob"
	"upload/shared/zombiekiller"

	blobopts "gocloud.dev/blob"
)

type CreateRequest struct {
	Metadata *File
	Rawdata  io.Reader
	Bucket   string
}

type CreateResponse struct {
	Metadata *File
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (CreateResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	metadata := make(chan error)
	go func() {
		metadata <- s.Repo.Insert(ctx, req.Metadata)
	}()

	rawdata := make(chan error)
	go func() {
		opts := &blobopts.WriterOptions{ContentType: req.Metadata.ContentType}
		rawdata <- s.Blob.Upload(ctx, req.Bucket, req.Metadata.ID.String(), req.Rawdata, opts)
	}()

	sentinel := struct {
		mu          sync.Mutex
		responded   bool
		errored     bool
		errResponse error
	}{}

	for i := 0; i < 2; i++ {
		select {
		case err := <-metadata:
			sentinel.mu.Lock()
			if err != nil {
				cancel()
				log.Printf("file: service: create: metadata: %v\n", err)
				if sentinel.responded && !sentinel.errored {
					s.Thrash <- zombiekiller.KillOperation{Killer: s.Blob, Target: &blob.Location{Bucket: req.Bucket, Key: req.Metadata.ID.String()}}
				}
				sentinel.errored = true
				sentinel.errResponse = err
			}
			sentinel.responded = true
			sentinel.mu.Unlock()
		case err := <-rawdata:
			sentinel.mu.Lock()
			if err != nil {
				cancel()
				log.Printf("file: service: create: rawdata: %v\n", err)
				if sentinel.responded && !sentinel.errored {
					s.Thrash <- zombiekiller.KillOperation{Killer: s.Repo, Target: req.Metadata.ID}
				}
				sentinel.errored = true
				sentinel.errResponse = err
			}
			sentinel.responded = true
			sentinel.mu.Unlock()
		}
	}

	if sentinel.errored {
		return CreateResponse{}, sentinel.errResponse
	}
	return CreateResponse{req.Metadata}, nil
}
