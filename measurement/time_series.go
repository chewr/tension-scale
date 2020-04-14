package measurement

import (
	"errors"
	"periph.io/x/periph/experimental/conn/analog"
	"sync"
	"time"
)

type TimeSeriesDevice struct {
	samples <-chan TimeSeriesSample
}

func NewTimeSeriesDevice(c <-chan TimeSeriesSample) (*TimeSeriesDevice) {
	return &TimeSeriesDevice{samples:c}
}

func (d *TimeSeriesDevice) Stream() <-chan TimeSeriesSample {
	return d.samples
}

func (d *TimeSeriesDevice) Read() (TimeSeriesSample, error) {
	sample, ok := <-d.samples
	if !ok {
		return sample, errors.New("device closed")
	}
	return sample, nil
}

type TimeSeriesSample struct {
	analog.Sample
	time.Time
}

const RING_BUFFER_SIZE = 32

type TimeSeriesBuffer struct {
	sync.Mutex
	samples    [RING_BUFFER_SIZE]TimeSeriesSample
	rptr, wptr int
	ready      bool
}

func (b *TimeSeriesBuffer) Write(s TimeSeriesSample) {
	b.Lock()
	defer b.Unlock()
	if b.wptr == b.rptr && b.ready {
		b.rptr = (b.rptr + 1) % RING_BUFFER_SIZE
	}
	b.samples[b.wptr] = s
	b.ready = true
	b.wptr = (b.wptr + 1) % RING_BUFFER_SIZE
}

func (b *TimeSeriesBuffer) Read() (TimeSeriesSample, bool) {
	b.Lock()
	defer b.Unlock()
	if !b.ready {
		return TimeSeriesSample{}, false
	}
	ts := b.samples[b.rptr]
	b.rptr = (b.rptr + 1) % RING_BUFFER_SIZE
	b.ready = b.rptr != b.wptr
	return ts, true
}

func (b *TimeSeriesBuffer) ReadAll() []TimeSeriesSample {
	b.Lock()
	defer b.Unlock()
	if !b.ready {
		return nil
	}
	sz := ((b.rptr + RING_BUFFER_SIZE) - b.wptr) % RING_BUFFER_SIZE
	dt := make([]TimeSeriesSample, sz)
	for i := range dt {
		idx := (b.rptr + i) % RING_BUFFER_SIZE
		dt[i] = b.samples[idx]
	}
	b.rptr = (b.rptr + sz) % RING_BUFFER_SIZE
	b.ready = false
	return dt
}
