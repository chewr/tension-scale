package isometric

import (
	"context"
	"time"

	"github.com/chewr/tension-scale/hx711"
	"github.com/chewr/tension-scale/led"
	"github.com/chewr/tension-scale/loadcell"
	"periph.io/x/periph/conn/physic"
)

type Workout interface {
	Run(ctx context.Context, display *led.TrafficLight, loadCell loadcell.Sensor) error
}

var _ Workout = new(restInterval)

type restInterval time.Duration

func (r restInterval) Run(ctx context.Context, display *led.TrafficLight, _ loadcell.Sensor) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(r))
	defer cancel()

	display.RedOn()
	defer display.RedOff()

	<-ctx.Done()
	return nil
}

func RestInterval(r time.Duration) Workout {
	return restInterval(r)
}

var _ Workout = &workInterval{}

type workInterval struct {
	threshold        physic.Force
	timeUnderTension time.Duration
}

func (w workInterval) Run(ctx context.Context, display *led.TrafficLight, loadCell loadcell.Sensor) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second+2*w.timeUnderTension)
	defer cancel()
	defer display.RedOff()
	defer display.GreenOff()
	defer display.YellowOff()

	underTension := false
	var startTime time.Time
	for {
		// Read force
		r, err := loadCell.Read(ctx)
		switch err {
		case nil: // continue processing
		case hx711.ErrBadRead:
			continue // drop a bad reading and continue
		default:
			return err
		}

		// State machine transitions
		if r.Force < w.threshold {
			underTension = false
			if err := display.GreenOff(); err != nil {
				return err
			}
			if err := display.YellowOn(); err != nil {
				return err
			}
			continue
		}
		switch {
		case underTension:
			if r.Time.Sub(startTime) >= w.timeUnderTension && r.Force >= w.threshold {
				return nil
			}
		default:
			underTension = r.Force >= w.threshold
			startTime = r.Time
			display.YellowOff()
			display.GreenOn()
		}
	}
}

func WorkInterval(t physic.Force, tut time.Duration) Workout {
	return &workInterval{
		threshold:        t,
		timeUnderTension: tut,
	}
}

var _ Workout = composite{}

type composite []Workout

func (c composite) Run(ctx context.Context, display *led.TrafficLight, loadCell loadcell.Sensor) error {
	for _, w := range c {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if err := w.Run(ctx, display, loadCell); err != nil {
			return err
		}
	}
	return nil
}

func Composite(w ...Workout) Workout {
	return composite(w)
}
