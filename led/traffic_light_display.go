package led

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/chewr/tension-scale/display"
	"periph.io/x/periph/conn/gpio"
)

type trafficLightDisplay struct {
	mu           sync.Mutex
	currentState display.State
	leds         *trafficLight
}

func NewTrafficLightDisplay(grn, ylw, red gpio.PinOut) (display.AutoRefreshingModel, error) {
	if err := grn.Out(gpio.Low); err != nil {
		return nil, err
	}
	if err := ylw.Out(gpio.Low); err != nil {
		return nil, err
	}
	if err := red.Out(gpio.Low); err != nil {
		return nil, err
	}
	return &trafficLightDisplay{
		leds: &trafficLight{
			green:  grn,
			yellow: ylw,
			red:    red,
		},
	}, nil
}

func (d *trafficLightDisplay) UpdateState(state display.State) error {
	if state == nil {
		return errors.New("Unrecognized state")
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	d.currentState = state
	return nil
}

const ledRefreshRate = 10 * time.Millisecond

func (d *trafficLightDisplay) Start(ctx context.Context) {
	t := time.NewTicker(ledRefreshRate)
	defer t.Stop()
	for range t.C {
		currentState := d.updateCurrentState()
		// TODO(rchew) context logging
		_ = d.displayState(currentState)
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

func (d *trafficLightDisplay) updateCurrentState() display.State {
	d.mu.Lock()
	defer d.mu.Unlock()
	if expiring, ok := d.currentState.ExpiringState(); ok {
		if expiring.Deadline().Before(time.Now()) {
			d.currentState = expiring.Fallback()
		}
	}
	return d.currentState
}

func (d *trafficLightDisplay) displayState(state display.State) error {
	c, err := colorFromState(state)
	if err != nil {
		return err
	}
	return d.leds.setColor(c)
}
