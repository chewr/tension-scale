package interval

import (
	"context"
	"fmt"
	"time"

	"github.com/chewr/tension-scale/display"
	"github.com/chewr/tension-scale/display/input"
	"github.com/chewr/tension-scale/display/state"
	"github.com/chewr/tension-scale/errutil"
	"github.com/chewr/tension-scale/hx711"
	"github.com/chewr/tension-scale/isometric"
	"github.com/chewr/tension-scale/loadcell"
	"periph.io/x/periph/conn/physic"
)

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
	defer errutil.SwallowF(func() error { return model.UpdateState(state.Halt()) })

	// Tare + setup
	tareDur := 2 * time.Second
	if err := model.UpdateState(state.Tare(tareDur)); err != nil {
		return err
	}
	done := time.After(tareDur)
	time.Sleep(time.Second)
	if err := loadCell.Tare(ctx, 20); err != nil {
		return err
	}
	<-done

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second+2*w.timeUnderTension)
	defer cancel()

	if err := model.UpdateState(state.WaitForInput(input.ForceRequired(w.threshold), input.None())); err != nil {
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

		// Start clock once threshold is passed
		if !underTension && r.Force > w.threshold {
			underTension = true
			startTime = r.Time
		}

		// Loop branch control
		// this is done before updating model state to avoid negative durations
		if underTension && time.Now().Sub(startTime) > w.timeUnderTension {
			return updater.Finish(isometric.Success)
		}

		// Update model state
		var currentState display.State
		if underTension {
			currentState = state.Work(
				input.ForceRequired(w.threshold),
				input.ForceReceived(r.Force),
				w.timeUnderTension-(time.Now().Sub(startTime)),
			)
		} else {
			currentState = state.WaitForInput(
				input.ForceRequired(w.threshold),
				input.ForceReceived(r.Force),
			)
		}

		if err := model.UpdateState(currentState); err != nil {
			return err
		}
	}
}

func WorkInterval(t physic.Force, tut time.Duration) isometric.Workout {
	return &workInterval{
		threshold:        t,
		timeUnderTension: tut,
	}
}
