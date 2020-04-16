package measurement

import (
	"context"
	"github.com/chewr/tension-scale/measurement"
	"periph.io/x/periph/conn"
)

type Sensor interface {
	conn.Resource
	// Reset resets the hardware. If it returns nil, the
	// sensor should be ready to read and an immediately
	// subsequent call to IsReady can be expected to
	// return true
	Reset(ctx context.Context) error

	// IsReady returns whether data is ready to be read
	IsReady() bool

	// Read blocks until a read is available, and then
	// reads out data
	Read(ctx context.Context) (measurement.TimeSeriesSample, error)

	// TryRead reads immediately if data available, otherwise returns error
	TryRead() (measurement.TimeSeriesSample, error)
}

type StreamingSensor interface {
	Sensor

	// ReadContinuous continuously reads as data is available,
	// sending results over the returned channel
	ReadContinuous() <-chan measurement.TimeSeriesSample
}
