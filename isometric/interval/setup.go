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
	defer errutil.SwallowF(func() error { return model.UpdateState(state.Halt()) })

	tareDur := 5 * time.Second
	done := time.After(tareDur)
	if err := model.UpdateState(state.Tare(time.Now().Add(tareDur))); err != nil {
		return err
	}
	time.Sleep(time.Second)
	if err := loadCell.Tare(ctx, 40); err != nil {
		return err
	}
	<-done

	risingEdgeInput := &input.DynamicEdgeInput{}
	if err := model.UpdateState(state.WaitForInput(input.RisingEdge(100*physic.Newton), risingEdgeInput)); err != nil {
		return err
	}
	for {
		fs, err := loadcell.TryReadIgnoreErrors(ctx, loadCell, hx711.ErrBadRead)
		if err != nil {
			return err
		}
		risingEdgeInput.Update(fs)
		if fs.Force >= 20*physic.PoundForce {
			break
		}
	}
	return nil
}
