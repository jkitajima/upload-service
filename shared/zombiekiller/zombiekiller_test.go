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
	zk := &zombieKillerMock{}
	tgt := &targetMock{}

	cases := map[string]struct {
		sendCount, retryCount int
		rw                    io.ReadWriter
		outMsg                []byte
		tgt                   fmt.Stringer
		tgtString             string
		zk                    ZombieKiller
		zkErr                 error
	}{
		"nil zombiekiller must return msg to output": {1, 4, &bytes.Buffer{}, []byte("zombiekiller: either killer or target is nil. ignoring received operation\n"), tgt, "gopher", nil, nil},
		"nil target must return msg to output":       {1, 4, &bytes.Buffer{}, []byte("zombiekiller: either killer or target is nil. ignoring received operation\n"), nil, "gopher", zk, nil},
		"basic zombie killing":                       {1, 4, &bytes.Buffer{}, []byte("zombiekiller: zombie data was found and killed\n"), tgt, "gopher", zk, nil},
		"basic retry":                                {1, 4, &bytes.Buffer{}, []byte("zombiekiller: retrying kill operation (retry count: 1)\n"), tgt, "gopher", zk, ErrInternal},
	}

	for key, testcase := range cases {
		t.Run(key, func(t *testing.T) {
			require := require.New(t)
			done := make(chan any)
			opsChan := make(chan KillOperation, testcase.sendCount)

			if tgt != nil {
				tgt.On("String").Return(testcase.tgtString)
			}
			if zk != nil {
				zk.On("KillZombie", tgt).Return(testcase.zkErr)
			}

			for range testcase.sendCount {
				opsChan <- KillOperation{testcase.zk, testcase.tgt}
			}

			go ListenForKillOperations(done, opsChan, uint8(testcase.retryCount), testcase.rw)

			time.Sleep(1 * time.Second)
			close(done)
			b := make([]byte, len(testcase.outMsg))
			testcase.rw.Read(b)
			require.Contains(string(testcase.outMsg), string(b))
		})
	}
}
