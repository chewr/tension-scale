package led

import (
	"errors"
	"sync"
	"time"

	"github.com/chewr/tension-scale/display"
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

func colorFromState(state display.State) (color, error) {
	var baseColor color
	switch state.GetType() {
	case display.Halt:
		baseColor = 0
	case display.Rest:
		baseColor = red
	case display.Work:
		baseColor = green
	case display.Tare:
		baseColor = yellow
	case display.Wait:
		baseColor = yellow
	default:
		return baseColor, errors.New("State not recognized")
	}

	color := baseColor
	if inputDependent, ok := state.InputDependentState(); ok {
		if !inputDependent.Satisfied() {
			color |= yellow
		}
	}

	// blink red when expiring
	if expiringState, ok := state.ExpiringState(); ok {
		ttl := time.Until(expiringState.Deadline())
		if ttl > 0 && ttl <= 3*time.Second {
			// for the final 3 seconds of an expiring interval, blink
			// 250ms on/750ms off
			if ttl%time.Second > 750*time.Millisecond {
				color |= red
			} else {
				color &= ^red
			}
		}
	}
	return color, nil
}
