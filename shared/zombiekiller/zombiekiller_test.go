package zombiekiller

import (
	"bytes"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type zombieKillerMock struct{ mock.Mock }

func (zk *zombieKillerMock) KillZombie(key fmt.Stringer) error {
	args := zk.Called(key)
	return args.Error(0)
}

type targetMock struct{ mock.Mock }

func (t *targetMock) String() string {
	args := t.Called()
	return args.String(0)
}

func TestListenForKillOperations(t *testing.T) {
	cases := map[string]struct {
		sendCount, retryCount int
		rw                    io.ReadWriter
		outMsg                []byte
		tgt                   fmt.Stringer
		tgtString             string
		zk                    ZombieKiller
		zkErr                 error
	}{
		"nil zombiekiller":     {1, 1, &bytes.Buffer{}, []byte("zombiekiller: killer is nil. ignoring received operation\n"), &targetMock{}, "gopher", nil, nil},
		"nil target":           {1, 1, &bytes.Buffer{}, []byte("zombiekiller: target is nil. ignoring received operation\n"), nil, "gopher", &zombieKillerMock{}, nil},
		"basic kill operation": {1, 0, &bytes.Buffer{}, []byte("zombiekiller: zombie data was found and killed\n"), &targetMock{}, "gopher", &zombieKillerMock{}, nil},
		"max retries":          {1, 5, &bytes.Buffer{}, []byte("zombiekiller: maximum retries reached. could not delete zombie data\n"), &targetMock{}, "gopher", &zombieKillerMock{}, ErrInternal},
	}

	for key, testcase := range cases {
		t.Run(key, func(t *testing.T) {
			require := require.New(t)
			done := make(chan any)
			opsChan := make(chan KillOperation, testcase.sendCount)

			if testcase.tgt != nil {
				testcase.tgt.(*targetMock).On("String").Return(testcase.tgtString)
			}
			if testcase.zk != nil {
				testcase.zk.(*zombieKillerMock).On("KillZombie", testcase.tgt).Return(testcase.zkErr)
			}

			for range testcase.sendCount {
				opsChan <- KillOperation{testcase.zk, testcase.tgt}
			}

			var rw bytes.Buffer
			go ListenForKillOperations(done, opsChan, uint8(testcase.retryCount), &rw)

			time.Sleep(1 * time.Second)
			close(done)

			b := make([]byte, len(testcase.outMsg))
			rw.Read(b)
			require.Equal(string(testcase.outMsg), string(b))
		})
	}
}
