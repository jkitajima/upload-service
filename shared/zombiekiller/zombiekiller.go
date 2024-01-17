package zombiekiller

import (
	"errors"
	"fmt"
	"log"
	"time"
)

type ZombieKiller interface {
	KillZombie(key fmt.Stringer) error
}

type KillOperation struct {
	Killer ZombieKiller
	Target fmt.Stringer
}

func ListenForKillOperations(doneChan <-chan any, ops <-chan KillOperation) {
	log.Println("Zombie Killer is active and listening for incoming operations.")

	for {
		select {
		case <-doneChan:
			log.Println("zombie killer received a signal to stop listening for incoming operations")
			return
		case op := <-ops:
			fmt.Println("op received")
			time.Sleep(5 * time.Second)
			var retry bool
			var count uint

			if op.Killer == nil || op.Target == nil {
				log.Println("zombie killer: received an empty killer or target")
				goto End
			}

			for ; count <= 4; count++ {
				if retry {
					log.Printf("zombie killer: retrying operation (count %d)\n", count)
				}

				if err := op.Killer.KillZombie(op.Target); err != nil {
					log.Println(err)

					switch err {
					case ErrInternal:
						retry = true
					}
				} else {
					break
				}
			}

			if count >= 4 {
				return
			}

			log.Println("zombie killer: zombie data was found and killed")
		End:
		}
	}
}

var (
	ErrInternal = errors.New("zombie killer: killer detected a target but failed to kill zombie data")
	ErrNotFound = errors.New("zombie killer: target was not found")
)
