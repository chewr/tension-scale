package interval

import (
	"context"
	"fmt"
	"time"

	"github.com/chewr/tension-scale/display"
	"github.com/chewr/tension-scale/display/state"
	"github.com/chewr/tension-scale/isometric"
	"github.com/chewr/tension-scale/loadcell"
)

type restInterval time.Duration

func (r restInterval) String() string {
	return fmt.Sprintf("rest-%v", time.Duration(r))
}

func (r restInterval) Run(ctx context.Context, model display.Model, _ loadcell.Sensor, _ isometric.WorkoutRecorder) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(r))
	defer cancel()

	if err := model.UpdateState(state.Rest()); err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func RestInterval(r time.Duration) isometric.Workout {
	return restInterval(r)
}
