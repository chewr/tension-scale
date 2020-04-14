package loadcell

import (
	"context"
	"errors"
	"github.com/chewr/tension-scale/measurement"
	"periph.io/x/periph/experimental/conn/analog"
	"periph.io/x/periph/experimental/devices/hx711"
	"sync"
	"time"
)

func New(dev *hx711.Dev) *loadCell {
	return &loadCell{
		hx711: dev,
	}
}

type loadCell struct {
	sync.Mutex
	hx711 *hx711.Dev

	tare int32
}

func (l *loadCell) MeasureTimeSeries(ctx context.Context) <-chan measurement.TimeSeriesSample {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	output := make(chan measurement.TimeSeriesSample)

	go func() {
		defer func() { close(output) }()

		// outputs are time-sensitive, so we don't want to
		// stale our data by waiting too long to write it
		// to the output channel.
		maxWriteWait := time.Millisecond
		writeTimer := time.NewTimer(maxWriteWait)
		defer writeTimer.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			timeout := time.Second
			if deadline, ok := ctx.Deadline(); ok {
				timeout = time.Until(deadline)
			}
			value, err := l.hx711.ReadTimeout(timeout)
			if err != nil {
				continue
			}
			tss := measurement.TimeSeriesSample{
				Sample: analog.Sample{Raw: value},
				Time:   time.Now(),
			}
			writeTimer.Reset(maxWriteWait)
			select {
			case output <- tss:
			case <-writeTimer.C:
			}
		}
	}()

	return output
}

// Tare measures a number of samples, then sets the tare weight of
// the scale to the average value measured. If samples equals 0,
// then Tare reads until the context times out.
func (l *loadCell) Tare(ctx context.Context, samples int) error {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	input := l.hx711.ReadContinuous()
	data := make([]int32, 0, samples)
read:
	for {
		select {
		case d, ok := <-input:
			if !ok {
				break read
			}
			data = append(data, d.Raw)
			if len(data) == samples {
				break read
			}
		case <-ctx.Done():
			if err := l.hx711.Halt(); err != nil {
				return err
			}
			break read
		}
	}

	if len(data) == 0 {
		return errors.New("Failed to collect tare data")
	}

	var total int32 = 0
	for _, d := range data {
		total += d
	}
	l.tare = total / int32(len(data))

	return nil
}

func (l *loadCell) setGain(gain int) error {
	switch gain {
	case 128:
		return l.hx711.SetInputMode(hx711.CHANNEL_A_GAIN_128)
	case 64:
		return l.hx711.SetInputMode(hx711.CHANNEL_A_GAIN_64)
	case 32:
		return l.hx711.SetInputMode(hx711.CHANNEL_B_GAIN_32)
	default:
		return errors.New("Invalid gain selected")
	}
}
