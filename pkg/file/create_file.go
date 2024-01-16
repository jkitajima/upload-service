package file

import (
	"context"
	"io"
	"log"

	"gocloud.dev/blob"
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

	metadata := make(chan error, 1)
	go func() {
		metadata <- s.Repo.Insert(ctx, req.Metadata)
	}()

	rawdata := make(chan error, 1)
	go func() {
		opts := &blob.WriterOptions{ContentType: req.Metadata.ContentType}
		rawdata <- s.Blob.Upload(ctx, req.Bucket, req.Metadata.ID.String(), req.Rawdata, opts)
	}()

	for i := 0; i < 2; i++ {
		select {
		case err := <-metadata:
			if err != nil {
				cancel()
				log.Printf("file: service: create: metadata: %v\n", err)

				if err := <-rawdata; err != nil {
					// both metadata and rawdata errored
					log.Printf("file: service: create: metadata: rawdata: %v\n", err)
				}

				// metadata ERRORED but RAWDATA dont
				// send ERR to ge queue

				return CreateResponse{}, err
			}
		case err := <-rawdata:
			if err != nil {
				cancel()
				log.Printf("file: service: create: rawdata: %v\n", err)

				if err := <-metadata; err != nil {
					// both metadata and rawdata errored
					log.Printf("file: service: create: rawdata: metadata: %v\n", err)
				}

				// metadata ERRORED but RAWDATA dont
				// send ERR to ge queue

				return CreateResponse{}, err
			}
		}
	}

	return CreateResponse{req.Metadata}, nil
}
