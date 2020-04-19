package led

import (
	"sync"

	"periph.io/x/periph/conn/gpio"
)

type TrafficLight struct {
	mu                 sync.Mutex
	Green, Yellow, Red gpio.PinOut
}

func NewTrafficLight(grn, ylw, red gpio.PinOut) (*TrafficLight, error) {
	if err := grn.Out(gpio.Low); err != nil {
		return nil, err
	}
	if err := ylw.Out(gpio.Low); err != nil {
		return nil, err
	}
	if err := red.Out(gpio.Low); err != nil {
		return nil, err
	}
	return &TrafficLight{
		Green:  grn,
		Yellow: ylw,
		Red:    red,
	}, nil
}

func (l *TrafficLight) GreenOn() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.Green.Out(gpio.High)
}

func (l *TrafficLight) GreenOff() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.Green.Out(gpio.Low)
}

func (l *TrafficLight) YellowOn() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.Yellow.Out(gpio.High)
}

func (l *TrafficLight) YellowOff() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.Yellow.Out(gpio.Low)
}

func (l *TrafficLight) RedOn() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.Red.Out(gpio.High)
}

func (l *TrafficLight) RedOff() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.Red.Out(gpio.Low)
}
