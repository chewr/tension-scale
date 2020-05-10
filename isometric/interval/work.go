package interval

import (
	"context"
	"fmt"
	"time"

	"github.com/chewr/tension-scale/display"
	"github.com/chewr/tension-scale/display/state"
	"github.com/chewr/tension-scale/hx711"
	"github.com/chewr/tension-scale/isometric"
	"github.com/chewr/tension-scale/loadcell"
	"periph.io/x/periph/conn/physic"
)

var _ isometric.Workout = &workInterval{}

type workInterval struct {
	threshold        physic.Force
	timeUnderTension time.Duration
}

func (w workInterval) String() string {
	return fmt.Sprintf("static-%v-%s",
		w.timeUnderTension,
		w.threshold.String(),
	)
}

func (w workInterval) Run(ctx context.Context, model display.Model, loadCell loadcell.Sensor, recorder isometric.WorkoutRecorder) error {
	defer model.UpdateState(state.Halt())

	// Tare + setup
	if err := model.UpdateState(state.Tare()); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	if err := loadCell.Tare(ctx, 20); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second+2*w.timeUnderTension)
	defer cancel()

	if err := model.UpdateState(state.WaitForInput()); err != nil {
		return err
	}

	updater, err := recorder.Start(ctx, w.String())
	if err != nil {
		return err
	}
	defer updater.Close()

	underTension := false
	var startTime time.Time
	for {
		// update model
		switch {
		case underTension:
			if err := model.UpdateState(state.Work()); err != nil {
				return err
			}
		default:
			if err := model.UpdateState(state.WaitForInput()); err != nil {
				return err
			}
		}

		// Read force
		r, err := loadCell.Read(ctx)
		switch err {
		case nil: // continue processing
		case hx711.ErrBadRead:
			continue // drop a bad reading and continue
		case context.DeadlineExceeded:
			return updater.Finish(isometric.Failure)
		default:
			return err
		}

		// record data
		if err := updater.Write(r); err != nil {
			return err
		}

		// State machine transitions
		if underTension &&
			// relax threshold to 75% once already under tension
			4*r.Force > 3*w.threshold {
			if r.Time.Sub(startTime) >= w.timeUnderTension && r.Force >= w.threshold {
				return updater.Finish(isometric.Success)
			}
		} else {
			underTension = r.Force > w.threshold
			if underTension { // will be true at first transition only
				startTime = r.Time
			}
		}
	}
}

func WorkInterval(t physic.Force, tut time.Duration) isometric.Workout {
	return &workInterval{
		threshold:        t,
		timeUnderTension: tut,
	}
}
