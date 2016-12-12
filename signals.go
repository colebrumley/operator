package operator

import (
	"os"
	"os/signal"
	"sync"

	log "github.com/Sirupsen/logrus"
)

// WatchForSignals monitors for os kill signals and runs arbitrary
// cleanup callbacks before exiting.
func WatchForSignals(callbacks ...func()) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)
	for {
		select {
		case s := <-signals:
			log.Warningln("Trapped signal:", s)
			if len(callbacks) > 0 {
				wg := sync.WaitGroup{}
				for _, c := range callbacks {
					wg.Add(1)
					go func(f func()) {
						f()
						wg.Done()
					}(c)
				}
				wg.Wait()
			}
			os.Exit(1)
		}
	}
}
