package loadcell

import (
	"context"
	"sync"

	"github.com/chewr/tension-scale/hx711"
)

type hx711Sensor struct {
	// immutable
	mu          sync.Mutex
	hx          hx711.V2
	calibration Calibration

	// mutable
	tare int64
}

func NewHx711(hx hx711.V2, calibration Calibration) Sensor {
	return &hx711Sensor{
		hx:          hx,
		calibration: calibration,
	}
}

func (s *hx711Sensor) Tare(ctx context.Context, samples int) error {
	if samples <= 0 {
		return ErrNotEnoughSamples
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	var total int64 = 0
	for i := 0; i < samples; i++ {
		r, err := s.hx.Read(ctx)
		if err != nil {
			return err
		}
		total += int64(r.Raw)
	}
	s.tare = total / int64(samples)
	return nil
}

func (s *hx711Sensor) Reset(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.hx.Reset(ctx); err != nil {
		return err
	}
	return nil
}

func (s *hx711Sensor) Halt() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.hx.Halt()
}

func (s *hx711Sensor) Read(ctx context.Context) (ForceSample, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, err := s.hx.Read(ctx)
	if err != nil {
		return ForceSample{}, err
	}
	return ForceSample{
		Force: s.calibration.ToForce(int64(r.Raw) - s.tare),
		Time:  r.Time,
	}, nil
}
