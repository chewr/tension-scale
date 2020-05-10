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

func SetupInterval(duration time.Duration) isometric.Workout {
	return setupInterval(duration)
}

type setupInterval time.Duration

func (s setupInterval) String() string {
	return fmt.Sprintf("setup-%v", time.Duration(s))
}

func (s setupInterval) Run(ctx context.Context, model display.Model, loadCell loadcell.Sensor, _ isometric.WorkoutRecorder) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s))
	defer cancel()
	defer model.UpdateState(state.Halt())

	if err := model.UpdateState(state.Tare()); err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	if err := loadCell.Tare(ctx, 40); err != nil {
		return err
	}

	if err := model.UpdateState(state.WaitForInput()); err != nil {
		return err
	}
	for {
		fs, err := loadcell.TryReadIgnoreErrors(ctx, loadCell, hx711.ErrBadRead)
		if err != nil {
			return err
		}
		if fs.Force >= 20*physic.PoundForce {
			break
		}
	}
	return nil
}
