package interval

import (
	"context"
	"strings"

	"github.com/chewr/tension-scale/display"
	"github.com/chewr/tension-scale/isometric"
	"github.com/chewr/tension-scale/loadcell"
)

var _ isometric.Workout = composite{}

type composite []isometric.Workout

func (c composite) String() string {
	workouts := []isometric.Workout(c)
	s := make([]string, len(workouts))
	for i, w := range workouts {
		s[i] = w.String()
	}
	return strings.Join(s, ",")
}

func (c composite) Run(ctx context.Context, model display.Model, loadCell loadcell.Sensor, recorder isometric.WorkoutRecorder) error {
	for _, w := range c {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if err := w.Run(ctx, model, loadCell, recorder); err != nil {
			return err
		}
	}
	return nil
}

func Composite(w ...isometric.Workout) isometric.Workout {
	return composite(w)
}
