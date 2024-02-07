package file

import (
	"context"
	"testing"

	"upload/shared/zombiekiller"

	"github.com/google/uuid"
)

func TestDeleteByID(t *testing.T) {
	ctx := context.Background()
	repo := new(repoerMock)
	stg := new(storagerMock)
	thrash := make(chan<- zombiekiller.KillOperation)
	s := &Service{repo, stg, thrash}

	id, _ := uuid.Parse("6bd1edba-bfe6-4a18-b2f3-db1dc4043b32")

	cases := map[string]struct {
		inReq  DeleteByIDRequest
		outErr error
	}{
		"basic": {DeleteByIDRequest{id, "bucketName"}, nil},
	}

	for key, testcase := range cases {
		t.Run(key, func(t *testing.T) {
			repo.On("DeleteByID", testcase.inReq.ID).Return(testcase.outErr)
			stg.On("Delete", testcase.inReq.Bucket, testcase.inReq.ID.String()).Return(testcase.outErr)

			if err := s.DeleteByID(ctx, testcase.inReq); err != testcase.outErr {
				t.Errorf("file: service: test_delete_by_id: error mismatched (result = %v, expected = %v)\n", err, testcase.outErr)
			}
		})
	}
}
