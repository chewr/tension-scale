package backcompat

import (
	"context"
	"github.com/chewr/tension-scale/hx711"
	"periph.io/x/periph/conn/pin"
	"periph.io/x/periph/experimental/conn/analog"
	periphimpl "periph.io/x/periph/experimental/devices/hx711"
	"time"
)

func HX711FromV2(v2 hx711.V2) HX711 {
	return &hx711Bridge{v2}
}

type hx711Bridge struct {
	v2 hx711.V2
}

func (b *hx711Bridge) ReadTimeout(timeout time.Duration) (int32, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()
	r, err := b.v2.Read(ctx)
	return r.Raw, err
}
func (b *hx711Bridge) SetInputMode(inputMode periphimpl.InputMode) error {
	var gain hx711.Gain
	switch inputMode {
	case periphimpl.CHANNEL_A_GAIN_128:
		gain = hx711.ChannelA128
	case periphimpl.CHANNEL_A_GAIN_64:
		gain = hx711.ChannelA64
	case periphimpl.CHANNEL_B_GAIN_32:
		gain = hx711.ChannelB32
	default:
		return hx711.ErrGainUnavailable
	}
	return b.v2.SetGain(gain)
}
func (b *hx711Bridge) IsReady() bool { return b.v2.IsReady() }

func (b *hx711Bridge) Range() (analog.Sample, analog.Sample) {
	return analog.Sample{Raw: -(1 << 23)}, analog.Sample{Raw: 1 << 23}
}

func (b *hx711Bridge) Read() (analog.Sample, error) {
	r, err := b.v2.TryRead()
	return r.Sample, err
}

func (b *hx711Bridge) Name() string { return b.v2.String() }
func (b *hx711Bridge) Number() int  { return -1 }
func (b *hx711Bridge) Function() string {
	return string(b.Func())
}

func (b *hx711Bridge) Func() pin.Func {
	return analog.ADC
}

func (b *hx711Bridge) String() string { return b.v2.String() }
func (b *hx711Bridge) ReadContinuous() <-chan analog.Sample {
	ch := b.v2.ReadContinuous()
	out := make(chan analog.Sample)
	go func() {
		defer close(out)
		for r := range ch {
			out <- r.Sample
		}
	}()
	return out
}
func (b *hx711Bridge) Halt() error { return b.v2.Halt() }
