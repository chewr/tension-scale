package maxhangs

import (
	"context"
	"github.com/chewr/tension-scale/led"
	"github.com/chewr/tension-scale/loadcell"
	"github.com/chewr/tension-scale/measurement"
	"time"
)

type Workout interface {
	Run(ctx context.Context, display led.TrafficLight, input measurement.TimeSeriesDevice) error
}

var _ Workout = new(restInterval)

type restInterval time.Duration

func (r restInterval) Run(ctx context.Context, display led.TrafficLight, _ measurement.TimeSeriesDevice) error {
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
	threshold        loadcell.ForceUnit
	timeUnderTension time.Duration
}

func (w workInterval) Run(ctx context.Context, display led.TrafficLight, input measurement.TimeSeriesDevice) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	defer display.GreenOff()
	defer display.YellowOff()

	underTension := false
	var startTime time.Time
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case sample := <-input.Stream():
			if sample.Add(time.Second / 2).Before(time.Now()) {
				// throw out old samples
				continue
			}
			if loadcell.ForceUnit(sample.Raw) >= w.threshold {
				if underTension && sample.Time.After(startTime.Add(w.timeUnderTension)) {
					return nil
				}
				if !underTension {
					underTension = true
					startTime = sample.Time
					display.YellowOff()
					display.GreenOn()
				}
			} else {
				underTension = false
				display.GreenOff()
				display.YellowOn()
			}
		}
	}
}

func WorkInterval(t loadcell.ForceUnit, tut time.Duration) Workout {
	return &workInterval{
		threshold:        t,
		timeUnderTension: tut,
	}
}

var _ Workout = composite{}

type composite []Workout

func (c composite) Run(ctx context.Context, display led.TrafficLight, input measurement.TimeSeriesDevice) error {
	for _, w := range c {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if err := w.Run(ctx, display, input); err != nil {
			return err
		}
	}
	return nil
}

func Composite(w ...Workout) Workout {
	return composite(w)
}
