package led

import (
	"time"
)

func Blink(on, off func() error, interval time.Duration) chan<- struct{} {
	t := time.NewTicker(interval)
	on()
	done := make(chan struct{})
	ticktock := 0
	go func() {
		for {
			select {
			case <-t.C:
				ticktock++
				if ticktock%2 == 1 {
					off()
				} else {
					on()
				}
			case <-done:
				t.Stop()
				return
			}
		}
	}()
	return done
}
