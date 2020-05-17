package led

import (
	"sync"

	"periph.io/x/periph/conn/gpio"
)

type color int

const (
	red = 1 << iota
	yellow
	green
)

type trafficLight struct {
	mu                 sync.Mutex
	green, yellow, red gpio.PinOut
}

func (l *trafficLight) setColor(c color) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if err := l.red.Out(c&red > 0); err != nil {
		return err
	}
	if err := l.green.Out(c&green > 0); err != nil {
		return err
	}
	if err := l.yellow.Out(c&yellow > 0); err != nil {
		return err
	}
	return nil
}
