package isometric

import (
	"context"
	"fmt"
	"time"

	"github.com/chewr/tension-scale/display"
	"github.com/chewr/tension-scale/display/state"
	"github.com/chewr/tension-scale/hx711"
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
func (t maxTest) Run(ctx context.Context, model display.Model, loadCell loadcell.Sensor, recorder WorkoutRecorder) error {
	defer model.UpdateState(state.Off)
	ctx, cancel := context.WithTimeout(ctx, time.Duration(t)*3)
	defer cancel()

	updater, err := recorder.Start(ctx, t.String())
	if err != nil {
		return nil
	}
	defer updater.Close()

	if err := model.UpdateState(state.Waiting); err != nil {
		return err
	}

	sw := &slidingWindow{
		dur: time.Duration(t),
	}
	trueMax := physic.Force(0)
	for {
		// Read force
		r, err := loadcell.TryReadIgnoreErrors(ctx, loadCell, hx711.ErrBadRead)
		if err != nil {
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
