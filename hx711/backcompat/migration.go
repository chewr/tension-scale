package backcompat

import (
	"context"
	"github.com/chewr/tension-scale/hx711"
	"github.com/chewr/tension-scale/measurement"
	"periph.io/x/periph/experimental/conn/analog"
	periphimpl "periph.io/x/periph/experimental/devices/hx711"
	"time"
)

func V2FromHX711(hx711 HX711) hx711.V2 {
	return &v2Bridge{hx711}
}

type v2Bridge struct {
	hx HX711
}

func (v *v2Bridge) String() string { return v.hx.Name() }

func (v *v2Bridge) Reset(ctx context.Context) error {
	for !v.hx.IsReady() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}
	return nil
}

// ReadContinuous implements hx711.V2
//
// This implementation prefers fresh data from the underlying
// HX711. If new data becomes available while it is blocking on
// sending old data, it will discard the old data. This is
// because the timestamp must be appended to the data as soon
// as possible after the data has been read
func (v *v2Bridge) ReadContinuous() <-chan measurement.TimeSeriesSample {
	ch := v.hx.ReadContinuous()
	ret := make(chan measurement.TimeSeriesSample)
	go func() {
		defer close(ret)
		var (
			in       analog.Sample
			out      measurement.TimeSeriesSample
			chanOpen = true
		)
		in, chanOpen = <-ch
		for chanOpen {
			out = measurement.TimeSeriesSample{
				Sample: in,
				Time:   time.Now(),
			}
			select {
			case in, chanOpen = <-ch:
			case ret <- out:
				in, chanOpen = <-ch
			}
		}
	}()
	return ret
}

func (v *v2Bridge) Halt() error { return v.hx.Halt() }

func (v *v2Bridge) IsReady() bool { return v.hx.IsReady() }

func (v *v2Bridge) Read(ctx context.Context) (measurement.TimeSeriesSample, error) {
	if deadline, ok := ctx.Deadline(); ok {
		s, err := v.hx.ReadTimeout(time.Until(deadline))
		return measurement.TimeSeriesSample{
			Sample: analog.Sample{Raw: s},
			Time:   time.Now(),
		}, err
	}
	if err := v.Reset(ctx); err != nil {
		return measurement.TimeSeriesSample{}, err
	}
	s, err := v.hx.Read()
	return measurement.TimeSeriesSample{
		Sample: s,
		Time:   time.Now(),
	}, err
}

func (v *v2Bridge) TryRead() (measurement.TimeSeriesSample, error) {
	if !v.hx.IsReady() {
		return measurement.TimeSeriesSample{}, hx711.ErrNotReady
	}
	s, err := v.hx.Read()
	return measurement.TimeSeriesSample{
		Sample: s,
		Time:   time.Now(),
	}, err
}

func (v *v2Bridge) SetGain(_ context.Context, g hx711.Gain) error {
	var inputMode periphimpl.InputMode
	switch g {
	case hx711.ChannelA128:
		inputMode = periphimpl.CHANNEL_A_GAIN_128
	case hx711.ChannelA64:
		inputMode = periphimpl.CHANNEL_A_GAIN_64
	case hx711.ChannelB32:
		inputMode = periphimpl.CHANNEL_B_GAIN_32
	default:
		return hx711.ErrGainUnavailable
	}
	return v.hx.SetInputMode(inputMode)
}

func (v *v2Bridge) Range() (analog.Sample, analog.Sample) {
	return analog.Sample{Raw: -(1 << 23)}, analog.Sample{Raw: 1 << 23}
}
