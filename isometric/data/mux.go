package data

import (
	"context"
	"sync"

	"github.com/chewr/tension-scale/isometric"
	"github.com/chewr/tension-scale/loadcell"
)

type multiplexingRecorder struct {
	mu        sync.Mutex
	recorders []isometric.WorkoutRecorder
}

func MultiRecorder(recorders ...isometric.WorkoutRecorder) isometric.WorkoutRecorder {
	return &multiplexingRecorder{
		recorders: recorders,
	}
}

func (r *multiplexingRecorder) Start(ctx context.Context, name string) (_ isometric.WorkoutUpdater, rErr error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	updaters := make([]isometric.WorkoutUpdater, 0, len(r.recorders))
	defer func() {
		if rErr != nil {
			for _, u := range updaters {
				u.Close()
			}
		}
	}()
	for _, rc := range r.recorders {
		u, err := rc.Start(ctx, name)
		if err != nil {
			return nil, err
		}
		updaters = append(updaters, u)
	}
	return &multiplexingUpdater{
		updaters: updaters,
	}, nil
}

type multiplexingUpdater struct {
	mu       sync.Mutex
	updaters []isometric.WorkoutUpdater
}

func (u *multiplexingUpdater) Write(samples ...loadcell.ForceSample) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	for _, ud := range u.updaters {
		if err := ud.Write(samples...); err != nil {
			return err
		}
	}
	return nil
}

func (u *multiplexingUpdater) Finish(outcome isometric.WorkoutOutcome) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	for _, ud := range u.updaters {
		if err := ud.Finish(outcome); err != nil {
			return err
		}
	}
	return nil
}

func (u *multiplexingUpdater) Close() {
	u.mu.Lock()
	defer u.mu.Unlock()
	for _, ud := range u.updaters {
		ud.Close()
	}
}
