package file

import (
	"context"
	"testing"
	"time"

	"upload/shared/zombiekiller"

	"github.com/google/uuid"
)

func TestUpdateByID(t *testing.T) {
	ctx := context.Background()
	repo := new(repoerMock)
	stg := new(storagerMock)
	thrash := make(chan<- zombiekiller.KillOperation)
	s := &Service{repo, stg, thrash}

	id, _ := uuid.Parse("6bd1edba-bfe6-4a18-b2f3-db1dc4043b32")
	submittedAt, _ := time.Parse(time.RFC3339, "2024-01-07T20:49:49.767+00:00")
	updatedAt, _ := time.Parse(time.RFC3339, "2024-01-07T20:49:50.535+00:00")
	uploadedAt, _ := time.Parse(time.RFC3339, "2024-01-07T20:49:50.535+00:00")
	f := &File{id, "007", "12662", "randimg.jpeg", ".jpeg", "image/jpeg", 478474, "https://kours.storage.blob.com/files/randimg.jpeg", 0, "Me and my family visiting Brazil during Summer 2017", submittedAt, updatedAt, uploadedAt}

	cases := map[string]struct {
		inReq   UpdateByIDRequest
		outResp UpdateByIDResponse
		outErr  error
	}{
		"basic": {UpdateByIDRequest{id, f}, UpdateByIDResponse{f}, nil},
	}

	for key, testcase := range cases {
		t.Run(key, func(t *testing.T) {
			repo.On("UpdateByID", testcase.inReq.ID, testcase.inReq.Metadata).Return(testcase.outErr)

			resp, err := s.UpdateByID(ctx, testcase.inReq)
			switch {
			case err != testcase.outErr:
				t.Errorf("file: service: test_update_by_id: error mismatched (result = %v, expected = %v)\n", err, testcase.outErr)
			case resp.Metadata != testcase.outResp.Metadata:
				t.Errorf("file: service: test_update_by_id: response mismatched (result = %v, expected = %v)\n", resp.Metadata, testcase.outResp.Metadata)
			}
		})
	}
}
