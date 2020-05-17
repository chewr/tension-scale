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

func NewTrafficLightDisplay(grn, ylw, red gpio.PinOut) (display.RunningModel, error) {
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

func (d *trafficLightDisplay) Run(ctx context.Context) {
	t := time.NewTicker(ledRefreshRate)
	defer t.Stop()
	for range t.C {
		// TODO(rchew) context logging
		d.refreshState()
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

func (d *trafficLightDisplay) refreshState() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	// TODO(rchew) blinking near expiry?
	// - don't want to deal with real time coding
	// - would make sense to apply "overlays" and deal with
	//   blinking at render time
	if expiring, ok := d.currentState.ExpiringState(); ok {
		if expiring.Deadline().Before(time.Now()) {
			d.currentState = expiring.Fallback()
		}
	}
	c, err := colorFromState(d.currentState)
	if err != nil {
		return err
	}
	if err := d.leds.setColor(c); err != nil {
		return err
	}
	return nil
}
