package loadcell

import (
	"context"
	"errors"
	"time"

	"periph.io/x/periph/conn/physic"
)

var (
	ErrNotEnoughSamples = errors.New("not enough samples to tare")
)

type ForceSample struct {
	physic.Force
	time.Time
}

// Sensor defines the API of a load cell sensor
type Sensor interface {
	Tare(ctx context.Context, samples int) error
	Reset(ctx context.Context) error
	Halt() error
	Read(ctx context.Context) (ForceSample, error)
}

func TryReadIgnoreErrors(ctx context.Context, sensor Sensor, ignore ...error) (ForceSample, error) {
	ignores := make(map[error]struct{})
	for _, e := range ignore {
		ignores[e] = struct{}{}
	}
	for {
		r, err := sensor.Read(ctx)
		if err != nil {
			if _, ok := ignores[err]; ok {
				continue
			}
		}
		return r, err
	}
}
