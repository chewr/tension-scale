package hx711

import (
	"context"
	"errors"
	"github.com/chewr/tension-scale/measurement"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/experimental/conn/analog"
	"periph.io/x/periph/host/cpu"
	"sync"
	"time"
)

const (
	// timing spec from hx711 datasheet:
	// https://components101.com/sites/default/files/component_datasheet/HX711%20Datasheet.pdf

	t1            = 100 * time.Nanosecond // T_1 typical time for data to be available after dout falling edge
	t2            = 100 * time.Nanosecond // T_2 typical time for dout to stabilize after pd_sck rising edge
	t3            = time.Microsecond      // T_3 typical pd_sck high time
	t4            = time.Microsecond      // T_4 typical pd_sck low time
	powerDownTime = 60 * time.Microsecond // time to hold pd_sck at HIGH to signal power down

	maxSampleRate = physic.Frequency(80 * physic.Hertz)
)

var (
	ErrNotImplemented = errors.New("not implemented")
	ErrStopped        = errors.New("hx711 is not powered on")
	ErrNotReady       = errors.New("data is not ready for reading")
)

// dev is a handle to a hx711.
type dev struct {
	// Immutable.
	name string
	clk  gpio.PinOut
	data gpio.PinIn

	// Mutable.
	mu        sync.Mutex
	inputMode Gain
	done      chan<- struct{}
	powerOn   bool
}

// New creates a new HX711 device.
//
// The data pin must support edge detection. If your pin doesn't natively
// support edge detection you can use PollEdge from gpioutil.
func New(clk gpio.PinOut, data gpio.PinIn) (V2, error) {
	if err := data.In(gpio.PullDown, gpio.FallingEdge); err != nil {
		return nil, err
	}
	if err := clk.Out(gpio.Low); err != nil {
		return nil, err
	}
	return &dev{
		name:      "hx711{" + clk.Name() + ", " + data.Name() + "}",
		inputMode: ChannelA128,
		clk:       clk,
		data:      data,
		done:      nil,
		powerOn:   true,
	}, nil
}

// ReadContinuous implements measurement.StreamingSensor
func (d *dev) ReadContinuous() <-chan measurement.TimeSeriesSample {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.done != nil {
		// read already in progress
		return nil
	}
	done := make(chan struct{})
	d.done = done

	if !d.powerOn {
		return nil
	}

	ctx, cancel := context.WithCancel(context.TODO())
	go func() {
		<-done
		cancel()
	}()

	out := make(chan measurement.TimeSeriesSample)
	go d.stream(ctx, out)
	return out
}

func (d *dev) stream(ctx context.Context, out chan<- measurement.TimeSeriesSample) {
	defer close(out)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		ts, err := d.Read(ctx)
		if err == nil {
			out <- ts
		}
	}
}

// Halt implements conn.Resource
func (d *dev) Halt() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.stop()
}

func (d *dev) stop() error {
	if d.done != nil {
		close(d.done)
		d.done = nil
	}
	if err := d.clk.Out(gpio.High); err != nil {
		return err
	}
	time.Sleep(powerDownTime)
	d.powerOn = false
	return nil
}

// IsReady implements measurement.Sensor
func (d *dev) IsReady() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.isReady()
}

func (d *dev) isReady() bool {
	return d.powerOn && d.data.Read() == gpio.Low
}

// Read implements measurement.Sensor
func (d *dev) Read(ctx context.Context) (measurement.TimeSeriesSample, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.blockingRead(ctx)
}

func (d *dev) blockingRead(ctx context.Context) (measurement.TimeSeriesSample, error) {
	if err := d.waitForReady(ctx); err != nil {
		return measurement.TimeSeriesSample{}, err
	}
	return d.readSample()
}

func (d *dev) waitForReady(ctx context.Context) error {
	// ensure device is turned on
	if !d.powerOn {
		return ErrStopped
	}
	// Coarse-grained wait; hx711 becomes ready 10-80Hz
	for !d.isReady() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		time.Sleep(maxSampleRate.Period() / 10)
	}

	// after DOUT falling edge, wait T_1 for data to be ready
	cpu.Nanospin(t1)
	return nil
}

func (d *dev) readSample() (measurement.TimeSeriesSample, error) {
	timestamp := time.Now()
	v, err := d.readRaw()
	if err != nil {
		return measurement.TimeSeriesSample{}, err
	}
	return measurement.TimeSeriesSample{
		Sample: analog.Sample{Raw: v},
		Time:   timestamp,
	}, ErrNotImplemented
}

func (d *dev) readRaw() (int32, error) {
	// Shift the 24-bit 2's compliment value.
	var value uint32
	for i := 0; i < 24; i++ {
		if err := d.clk.Out(gpio.High); err != nil {
			return 0, err
		}

		cpu.Nanospin(t2)
		level := d.data.Read()

		cpu.Nanospin(t3 - t2)
		if err := d.clk.Out(gpio.Low); err != nil {
			return 0, err
		}
		cpu.Nanospin(t4)

		value <<= 1
		if level {
			value |= 1
		}
	}

	// Pulse the clock 1-3 more times to set the new ADC mode.
	for i := 0; i < int(d.inputMode); i++ {
		if err := d.clk.Out(gpio.High); err != nil {
			return 0, err
		}
		cpu.Nanospin(t3)
		if err := d.clk.Out(gpio.Low); err != nil {
			return 0, err
		}
		cpu.Nanospin(t4)
	}
	// Convert the 24-bit 2's compliment value to a 32-bit signed value.
	return int32(value<<8) >> 8, nil
}

// TryRead implements measurement.Sensor
func (d *dev) TryRead() (measurement.TimeSeriesSample, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.powerOn {
		return measurement.TimeSeriesSample{}, ErrStopped
	}

	if !d.isReady() {
		return measurement.TimeSeriesSample{}, ErrNotReady
	}
	cpu.Nanospin(t1)
	return d.readSample()
}

// SetGain implements V2
func (d *dev) SetGain(ctx context.Context, g Gain) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	switch g {
	default:
		return ErrGainUnavailable
	case ChannelA128:
	case ChannelA64:
	case ChannelB32:
	}

	d.inputMode = g

	if !d.powerOn {
		return nil
	}

	// read and throw away one data point to set gain
	_, err := d.blockingRead(ctx)
	return err
}

// Range implements measurement.Sensor
func (d *dev) Range() (analog.Sample, analog.Sample) {
	return analog.Sample{Raw: -(1 << 23)}, analog.Sample{Raw: 1 << 23}
}

// String implements conn.Resource
func (d *dev) String() string { return d.name }

// Reset implements measurement.Sensor
func (d *dev) Reset(ctx context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if !d.powerOn {
		d.powerOn = true
		if err := d.clk.Out(gpio.Low); err != nil {
			return err
		}
	}

	if d.inputMode != ChannelA128 {
		if _, err := d.blockingRead(ctx); err != nil {
			return err
		}
	}

	return d.waitForReady(ctx)
}
