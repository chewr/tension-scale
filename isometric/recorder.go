package isometric

import (
	"context"
	"errors"

	"github.com/chewr/tension-scale/loadcell"
)

var (
	ErrWriteAfterClosed = errors.New("workout has already been closed")
	ErrNoData           = errors.New("no data to write out")
)

type WorkoutRecorder interface {
	Start(ctx context.Context, descriptor string) (WorkoutUpdater, error)
}

type WorkoutUpdater interface {
	Write(sample ...loadcell.ForceSample) error
	Finish(outcome WorkoutOutcome) error
	Close()
}
