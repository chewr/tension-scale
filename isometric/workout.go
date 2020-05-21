package isometric

import (
	"context"
	"fmt"

	"github.com/chewr/tension-scale/display"
	"github.com/chewr/tension-scale/loadcell"
)

type WorkoutOutcome string

const (
	Success WorkoutOutcome = "success"
	Pass    WorkoutOutcome = "pass"
	Failure WorkoutOutcome = "failure"
)

type Workout interface {
	fmt.Stringer
	Run(ctx context.Context, model display.Model, loadCell loadcell.Sensor, recorder WorkoutRecorder) error
}
