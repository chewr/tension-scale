package led

import (
	"errors"
	"periph.io/x/periph/conn/gpio"
	"sync"

	"github.com/chewr/tension-scale/display"
)

type trafficLightDisplay struct {
	mu   sync.Mutex
	leds *trafficLight
}

func NewTrafficLightDisplay(grn, ylw, red gpio.PinOut) (display.Model, error) {
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
	return errors.New("Not implemented")
}
