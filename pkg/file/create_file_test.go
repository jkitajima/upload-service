package file

import (
	"context"
	"strings"
	"testing"
	"time"

	"upload/shared/zombiekiller"

	"github.com/google/uuid"
	"gocloud.dev/blob"
)

func TestCreate(t *testing.T) {
	ctx := context.Background()
	repo := new(repoerMock)
	stg := new(storagerMock)
	thrash := make(chan<- zombiekiller.KillOperation)
	s := &Service{repo, stg, thrash}

	id, _ := uuid.Parse("6bd1edba-bfe6-4a18-b2f3-db1dc4043b32")
	submittedAt, _ := time.Parse(time.RFC3339, "2024-01-07T20:49:49.767+00:00")
	updatedAt, _ := time.Parse(time.RFC3339, "2024-01-07T20:49:50.535+00:00")
	uploadedAt, _ := time.Parse(time.RFC3339, "2024-01-07T20:49:50.535+00:00")
	f := &File{id, "17332", "96", "sunny_day.jpeg", ".jpeg", "image/jpeg", 32137, "https://kours.storage.blob.com/files/sunny_day.jpeg", 0, "Me and my family visiting Brazil during Summer 2017", submittedAt, updatedAt, uploadedAt}

	cases := map[string]struct {
		inReq   CreateRequest
		outResp CreateResponse
		outErr  error
	}{
		"basic": {CreateRequest{f, strings.NewReader("req1 test data"), "company"}, CreateResponse{f}, nil},
	}

	for key, testcase := range cases {
		t.Run(key, func(t *testing.T) {
			repo.On("Insert", testcase.inReq.Metadata).Return(testcase.outErr)
			stg.On("Upload", testcase.inReq.Bucket, testcase.inReq.Metadata.ID.String(), testcase.inReq.Rawdata, &blob.WriterOptions{ContentType: testcase.inReq.Metadata.ContentType}).Return(testcase.outErr)

			resp, err := s.Create(ctx, testcase.inReq)
			if err != testcase.outErr || resp != testcase.outResp {
				t.Errorf("deu ruim")
			}
		})
	}
}
