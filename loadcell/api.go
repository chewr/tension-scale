package loadcell

import (
	"context"
	"errors"
	"github.com/chewr/tension-scale/hx711"
	"github.com/chewr/tension-scale/measurement"
	"sync"
)

var (
	ErrNotEnoughSamples = errors.New("not enough samples to tare")
)

// Sensor defines the API of a load cell sensor
type Sensor interface {
	Tare(ctx context.Context, samples int) error
	Reset(ctx context.Context) error
	Halt() error
	Read(ctx context.Context) (measurement.TimeSeriesSample, error)
}

type hx711Sensor struct {
	sync.Mutex
	tare int32
	hx   hx711.V2
}

func (s *hx711Sensor) Tare(ctx context.Context, samples int) error {
	if samples <= 0 {
		return ErrNotEnoughSamples
	}
	s.Lock()
	defer s.Mutex.Unlock()
	var total int32 = 0
	for i := 0; i < samples; i++ {
		r, err := s.hx.Read(ctx)
		if err != nil {
			return err
		}
		total += r.Raw
	}
	s.tare = total / int32(samples)
	return nil
}

func (s *hx711Sensor) Reset(ctx context.Context) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	if err := s.hx.Reset(ctx); err != nil {
		return err
	}
	s.tare = 0
	return nil
}

func (s *hx711Sensor) Halt() error {
	return s.hx.Halt()
}

func (s *hx711Sensor) Read(ctx context.Context) (measurement.TimeSeriesSample, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	r, err := s.hx.Read(ctx)
	if err != nil {
		return measurement.TimeSeriesSample{}, err
	}
	r.Raw -= s.tare
	return r, err
}
