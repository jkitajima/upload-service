package composer

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
)

type serverMock struct {
	mock.Mock
}

func (s *serverMock) Mux() http.Handler {
	args := s.Called()
	return args.Get(0).(http.Handler)
}

func (s *serverMock) Prefix() string {
	args := s.Called()
	return args.String(0)
}

func TestCompose(t *testing.T) {
	s1, s2 := &serverMock{}, &serverMock{}
	servers := []Server{s1, s2}

	cases := map[string]struct {
		c       *Composer
		servers []Server
		mux     http.Handler
		err     error
	}{
		"basic compose": {NewComposer(), servers, http.NewServeMux(), nil},
	}

	for key, testcase := range cases {
		t.Run(key, func(t *testing.T) {
			for _, s := range testcase.servers {
				bytes := make([]byte, 4)
				rand.Read(bytes)
				tmpname := hex.EncodeToString(bytes)

				s.(*serverMock).On("Mux").Return(testcase.mux)
				s.(*serverMock).On("Prefix").Return("/" + tmpname)
			}

			if err := testcase.c.Compose(testcase.servers...); err != testcase.err {
				t.Errorf("composer: err mismatched (result = %v, expected = %v)\n", err, testcase.err)
			}
		})
	}
}
