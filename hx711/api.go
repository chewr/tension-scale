package hx711

import (
	"context"
	"errors"

	"github.com/chewr/tension-scale/measurement"
)

type Gain int

const (
	ChannelA128 Gain = 1
	ChannelA64  Gain = 3
	ChannelB32  Gain = 2
)

var (
	ErrGainUnavailable = errors.New("specified gain value is unavailable")
)

type V2 interface {
	measurement.StreamingSensor
	SetGain(ctx context.Context, gain Gain) error
}
