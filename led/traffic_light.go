package led

import (
	"periph.io/x/periph/conn/gpio"
)

type TrafficLight struct {
	Green, Yellow, Red gpio.PinOut
}

func (l TrafficLight) GreenOn() error {
	return l.Green.Out(gpio.High)
}

func (l TrafficLight) GreenOff() error {
	return l.Green.Out(gpio.Low)
}

func (l TrafficLight) YellowOn() error {
	return l.Yellow.Out(gpio.High)
}

func (l TrafficLight) YellowOff() error {
	return l.Yellow.Out(gpio.Low)
}

func (l TrafficLight) RedOn() error {
	return l.Red.Out(gpio.High)
}

func (l TrafficLight) RedOff() error {
	return l.Red.Out(gpio.Low)
}
