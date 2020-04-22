package isometric

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chewr/tension-scale/hx711"
	"github.com/chewr/tension-scale/led"
	"github.com/chewr/tension-scale/loadcell"
	"periph.io/x/periph/conn/physic"
)

type WorkoutOutcome string

const (
	Success WorkoutOutcome = "success"
	Pass    WorkoutOutcome = "pass"
	Failure WorkoutOutcome = "failure"
)

type Workout interface {
	fmt.Stringer
	Run(ctx context.Context, display *led.TrafficLight, loadCell loadcell.Sensor, recorder WorkoutRecorder) error
}

var _ Workout = new(restInterval)

type restInterval time.Duration

func (r restInterval) String() string {
	return fmt.Sprintf("rest-%v", time.Duration(r))
}

func (r restInterval) Run(ctx context.Context, display *led.TrafficLight, _ loadcell.Sensor, _ WorkoutRecorder) error {
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

func (w workInterval) String() string {
	return fmt.Sprintf("static-%v-%dlbs",
		w.timeUnderTension,
		w.threshold/physic.PoundForce,
	)
}

func (w workInterval) Run(ctx context.Context, display *led.TrafficLight, loadCell loadcell.Sensor, recorder WorkoutRecorder) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second+2*w.timeUnderTension)
	defer cancel()
	defer display.RedOff()
	defer display.GreenOff()
	defer display.YellowOff()

	updater, err := recorder.Start(ctx, w.String())
	defer updater.Close()
	if err != nil {
		return err
	}

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

		// record data
		if err := updater.Write(r); err != nil {
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
				return updater.Finish(Success)
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

func SetupInterval(d time.Duration) Workout {
	return setupInterval(d)
}

type setupInterval time.Duration

func (s setupInterval) String() string {
	return fmt.Sprintf("setup-%v", time.Duration(s))
}

func (s setupInterval) Run(ctx context.Context, display *led.TrafficLight, loadCell loadcell.Sensor, _ WorkoutRecorder) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s))
	defer cancel()

	go func() {
		done := led.Blink(display.YellowOn, display.YellowOff, time.Second/4)
		defer close(done)
		<-ctx.Done()
	}()

	for {
		fs, err := loadCell.Read(ctx)
		switch err {
		case nil: // continue processing
		case hx711.ErrBadRead:
			continue // drop a bad reading and continue
		default:
			return err
		}
		if fs.Force >= 20*physic.PoundForce {
			break
		}
	}
	return nil
}

var _ Workout = composite{}

type composite []Workout

func (c composite) String() string {
	workouts := []Workout(c)
	s := make([]string, len(workouts))
	for i, w := range workouts {
		s[i] = w.String()
	}
	return strings.Join(s, ",")
}

func (c composite) Run(ctx context.Context, display *led.TrafficLight, loadCell loadcell.Sensor, recorder WorkoutRecorder) error {
	for _, w := range c {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if err := w.Run(ctx, display, loadCell, recorder); err != nil {
			return err
		}
	}
	return nil
}

func Composite(w ...Workout) Workout {
	return composite(w)
}
