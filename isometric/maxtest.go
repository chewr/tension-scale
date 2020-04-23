package isometric

import (
	"context"
	"fmt"
	"time"

	"github.com/chewr/tension-scale/hx711"
	"github.com/chewr/tension-scale/led"
	"github.com/chewr/tension-scale/loadcell"
	"periph.io/x/periph/conn/physic"
)

func MaxTest(hold time.Duration) Workout {
	return maxTest(hold)
}

type maxTest time.Duration

func (t maxTest) String() string {
	return fmt.Sprintf("max-test-%v", time.Duration(t))
}
func (t maxTest) Run(ctx context.Context, display *led.TrafficLight, loadCell loadcell.Sensor, recorder WorkoutRecorder) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(t)*3+countdownSequenceTime)
	defer cancel()

	updater, err := recorder.Start(ctx, t.String())
	if err != nil {
		return nil
	}
	defer updater.Close()

	// display a little countdown sequence for user
	if err := countdownSequence(ctx, display); err != nil {
		return err
	}

	// lights on!
	if err := display.GreenOn(); err != nil {
		return err
	}
	defer display.GreenOff()
	defer display.YellowOff()
	defer display.RedOff()

	sw := &slidingWindow{
		dur: time.Duration(t),
	}
	trueMax := physic.Force(0)
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
		if r.Force > trueMax {
			trueMax = r.Force
		}
		sw.update(r)
		if sw.ready() {
			if trueMax > sw.maxForce() {
				return updater.Finish(Success)
			}
		}
	}
}

type slidingWindow struct {
	samples  []loadcell.ForceSample
	dur      time.Duration
	startPtr int
}

func (w *slidingWindow) update(sample loadcell.ForceSample) {
	w.samples = append(w.samples, sample)
	for ; sample.Time.Sub(w.samples[w.startPtr].Time) >= w.dur; w.startPtr++ {
	}
}

func (w *slidingWindow) ready() bool {
	return w.startPtr > 0
}

func (w *slidingWindow) maxForce() physic.Force {
	m := physic.Force(0)
	for i := w.startPtr; i < len(w.samples); i++ {
		if w.samples[i].Force > m {
			m = w.samples[i].Force
		}
	}
	return m
}

const countdownSequenceTime = time.Second * 5

func countdownSequence(ctx context.Context, display *led.TrafficLight) error {
	// setup - shut everything down
	if err := display.RedOff(); err != nil {
		return err
	}
	if err := display.YellowOff(); err != nil {
		return err
	}
	if err := display.GreenOff(); err != nil {
		return err
	}
	defer display.GreenOff()
	defer display.YellowOff()
	defer display.RedOff()

	// hold red for two seconds
	if err := display.RedOn(); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)

	// blink yellow slowly for two seconds
	for i := 0; i < 3; i++ {
		if err := display.YellowOn(); err != nil {
			return err
		}
		time.Sleep(125 * time.Millisecond)
		if err := display.YellowOff(); err != nil {
			return err
		}
		time.Sleep(375 * time.Millisecond)
	}

	// blink yellow quickly for one second
	for i := 0; i < 3; i++ {
		if err := display.YellowOn(); err != nil {
			return err
		}
		time.Sleep(125 * time.Millisecond)
		if err := display.YellowOff(); err != nil {
			return err
		}
		time.Sleep(125 * time.Millisecond)
	}

	// for the last second, hold yellow only
	if err := display.RedOff(); err != nil {
		return err
	}
	ch := led.Blink(display.YellowOn, display.YellowOff, time.Second/16)
	time.Sleep(time.Second)
	close(ch)
	return nil
}
