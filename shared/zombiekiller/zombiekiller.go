package zombiekiller

import (
	"errors"
	"fmt"
	"io"
)

type ZombieKiller interface {
	KillZombie(key fmt.Stringer) error
}

type KillOperation struct {
	Killer ZombieKiller
	Target fmt.Stringer
}

func ListenForKillOperations(doneChan <-chan any, ops <-chan KillOperation, retryCount uint8, out io.Writer) {
	for {
		select {
		case <-doneChan:
			return
		case op := <-ops:
			var retry bool
			var count uint8

			if op.Killer == nil || op.Target == nil {
				out.Write([]byte("zombiekiller: either killer or target is nil. ignoring received operation\n"))
				goto End
			}

			for ; count < retryCount; count++ {
				if retry {
					out.Write([]byte(fmt.Sprintf("zombiekiller: retrying kill operation (retry count: %d)\n", count)))
				}

				if err := op.Killer.KillZombie(op.Target); err != nil {
					switch err {
					case ErrInternal:
						retry = true
					}
				} else {
					break
				}
			}

			if count > retryCount {
				out.Write([]byte("zombiekiller: maximum retry count passed. aborting retries\n"))
				return
			}

			out.Write([]byte("zombiekiller: zombie data was found and killed\n"))
		End:
		}
	}
}

var (
	ErrInternal = errors.New("zombiekiller: killer detected a target but failed to kill zombie data")
	ErrNotFound = errors.New("zombiekiller: target was not found")
)
