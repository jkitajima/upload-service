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
			switch {
			case op.Killer == nil:
				fmt.Fprintln(out, "zombiekiller: killer is nil. ignoring received operation")
				goto End
			case op.Target == nil:
				fmt.Fprintln(out, "zombiekiller: target is nil. ignoring received operation")
				goto End
			}

			for i := 0; i <= int(retryCount); i++ {
				if i == int(retryCount) && retryCount > 0 {
					fmt.Fprintln(out, "zombiekiller: maximum retries reached. could not delete zombie data")
					break
				}

				if err := op.Killer.KillZombie(op.Target); err != nil {
					if err == ErrNotFound {
						break
					}
				} else {
					fmt.Fprintln(out, "zombiekiller: zombie data was found and killed")
					break
				}
			}

		End:
		}
	}
}

var (
	ErrInternal = errors.New("zombiekiller: killer detected a target but failed to kill zombie data")
	ErrNotFound = errors.New("zombiekiller: target was not found")
)
