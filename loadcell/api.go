package loadcell

import (
	"context"
	"errors"
	"periph.io/x/periph/conn/physic"
	"time"
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
