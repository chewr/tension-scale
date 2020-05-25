package led

import (
	"context"
	"sync"
	"time"

	"github.com/chewr/tension-scale/display"
	"github.com/chewr/tension-scale/display/stateimpl"
	"periph.io/x/periph/conn/gpio"
)

type trafficLightDisplay struct {
	mu     sync.Mutex
	ticker *time.Ticker

	stateimpl.StateHolder
	leds *trafficLight
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

const ledRefreshRate = 10 * time.Millisecond

func (d *trafficLightDisplay) Start(ctx context.Context) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.ticker != nil {
		return
	}
	d.ticker = time.NewTicker(ledRefreshRate)
	go d.run(ctx, d.ticker.C)
}

func (d *trafficLightDisplay) run(ctx context.Context, c <-chan time.Time) {
	defer d.stop()
	for range c {
		// TODO(rchew) context logging
		currentState, _ := d.GetCurrentState()
		_ = d.displayState(currentState)
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

func (d *trafficLightDisplay) stop() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.ticker != nil {
		d.ticker.Stop()
		d.ticker = nil
	}
}

func (d *trafficLightDisplay) displayState(state display.State) error {
	c, err := colorFromState(state)
	if err != nil {
		return err
	}
	return d.leds.setColor(c)
}
