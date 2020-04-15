package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/chewr/tension-scale/hx711"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"periph.io/x/periph/experimental/conn/analog"
	periphimpl "periph.io/x/periph/experimental/devices/hx711"
	"periph.io/x/periph/host/rpi"
	"sync"
	"time"
)

var (
	readCmd = &cobra.Command{
		Use:   "read",
		Short: "Read data in from the hx711",
		RunE:  doRead,
	}
)

var (
	ErrNotImplemented = errors.New("command not implemented")
	ErrStopped        = errors.New("driver not running")
)

var (
	PIN_HX711_SCLK = rpi.P1_31 // GPIO6
	PIN_HX711_DOUT = rpi.P1_29 // GPIO5
)

func loadHx711(cmd *cobra.Command) (hx711.HX711, error) {
	if viper.GetBool(flagUsePeriphImplementation) {
		cmd.Println("Using periph hx711 driver implementation")
		return periphimpl.New(PIN_HX711_SCLK, PIN_HX711_DOUT)
	}
	cmd.Println("Using custom hx711 driver implementation")
	hx, err := hx711.New(PIN_HX711_SCLK, PIN_HX711_DOUT)
	if err != nil {
		return nil, err
	}
	if viper.GetBool(flagReset) {
		cmd.Println("Resetting the hx711 module")
		err = hx.Reset()
	}
	return hx, err
}

func initHx711(cmd *cobra.Command) (hx711.HX711, error) {
	hx, err := loadHx711(cmd)
	if err != nil {
		return nil, err
	}
	switch viper.GetInt(flagGain) {
	case 32:
		if err := hx.SetInputMode(periphimpl.CHANNEL_B_GAIN_32); err != nil {
			return nil, err
		}
	case 64:
		if err := hx.SetInputMode(periphimpl.CHANNEL_A_GAIN_64); err != nil {
			return nil, err
		}
	default: // 128 gain is configured by default
	}
	return hx, nil
}

func doRead(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()
	validateFlags(cmd)

	sp, err := getSampleProducer(cmd)
	if err != nil {
		return err
	}
	sp.Start(ctx)
	defer sp.Stop()

	return consume(ctx, cmd, sp)
}

func getSampleProducer(cmd *cobra.Command) (SampleProducer, error) {
	hx, err := initHx711(cmd)
	if err != nil {
		cmd.PrintErrln(err)
		return nil, err
	}

	if viper.GetBool(flagContinuous) {
		cmd.Println("Creating continuous sample producer")
		return NewContinuousSampleProducer(hx), nil
	} else if viper.GetBool(flagInstantaneousRead) {
		cmd.Println("Creating on-demand sample producer (instantaneous)")
		return NewOnDemandSampleProvider(hx, false), nil
	}
	cmd.Println("Creating on-demand sample producer (patient)")
	return NewOnDemandSampleProvider(hx, true), nil
}

type sample struct {
	value     int32
	timestamp time.Time
}

type SampleProducer interface {
	Start(ctx context.Context)
	Read(ctx context.Context) (sample, error)
	Stop() error
}

var _ SampleProducer = &continuousSampleProducer{}

func NewContinuousSampleProducer(hx hx711.HX711) SampleProducer {
	return &continuousSampleProducer{
		hx: hx,
	}
}

type continuousSampleProducer struct {
	sync.Mutex
	hx hx711.HX711
	ch <-chan sample
}

func (p *continuousSampleProducer) Start(ctx context.Context) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	if p.ch != nil {
		return
	}
	// buffer the channel so data doesn't go stale while writes block
	inCh := p.hx.ReadContinuous()
	outCh := make(chan sample, 1000)

	go pipe(ctx, inCh, outCh)

	p.ch = outCh
	go func() {
		<-ctx.Done()
		p.Stop()
	}()
}

func (p *continuousSampleProducer) Stop() error {
	err := p.hx.Halt()
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	p.ch = nil
	return err
}

func pipe(ctx context.Context, inCh <-chan analog.Sample, outCh chan<- sample) {
	defer close(outCh)
	for v := range inCh {
		s := sample{
			value:     v.Raw,
			timestamp: time.Now(),
		}
		select {
		case outCh <- s:
		case <-ctx.Done():
			return
		}
	}
}

func (p *continuousSampleProducer) Read(ctx context.Context) (sample, error) {
	ch := p.getChannel()
	for ch == nil {
		select {
		case <-ctx.Done():
			return sample{}, ctx.Err()
		default:
		}
		ch = p.getChannel()
	}
	select {
	case <-ctx.Done():
		return sample{}, ctx.Err()
	case s, ok := <-ch:
		if ok {
			return s, nil
		}
		return s, ErrStopped
	}
}

func (p *continuousSampleProducer) getChannel() <-chan sample {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	return p.ch
}

var _ SampleProducer = &onDemandSampleProducer{}

func NewOnDemandSampleProvider(hx hx711.HX711, useTimeout bool) SampleProducer {
	return &onDemandSampleProducer{
		hx:         hx,
		useTimeout: useTimeout,
	}
}

type onDemandSampleProducer struct {
	sync.Mutex
	stopped    bool
	hx         hx711.HX711
	useTimeout bool
}

func (p *onDemandSampleProducer) Start(_ context.Context) {}

func (p *onDemandSampleProducer) Stop() error {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	p.stopped = true
	return nil
}

func (p *onDemandSampleProducer) Read(ctx context.Context) (sample, error) {
	p.Lock()
	defer p.Mutex.Unlock()
	if p.stopped {
		return sample{}, ErrStopped
	}
	var (
		s   sample
		err error
	)
	if p.useTimeout {
		var value int32
		value, err = p.hx.ReadTimeout(time.Second)
		s = sample{
			value:     value,
			timestamp: time.Now(),
		}
	} else {
		var as analog.Sample
		as, err = p.hx.Read()
		s = sample{
			value:     as.Raw,
			timestamp: time.Now(),
		}
	}
	return s, err
}

func consume(ctx context.Context, cmd *cobra.Command, p SampleProducer) error {
	for {
		s, err := p.Read(ctx)
		switch err {
		case nil:
			break
		case ErrStopped:
			cmd.Println("Device stopped")
			return nil
		default:
			cmd.PrintErrln("failed to read sample: ", err)
			return err
		}

		cmd.Println(printSample(s))

		select {
		case <-ctx.Done():
			return nil
		}
	}
}

func printSample(s sample) string {
	// now: reading | 24bit | timestamp | delta
	now := time.Now()
	const timeFormat = "15:04:05.000000"
	return fmt.Sprintf("%s: %10d | %s | %s | %14d",
		now.Format(timeFormat),
		s.value,
		to24BitBinary(s.value),
		s.timestamp.Format(timeFormat),
		now.Sub(s.timestamp))
}

func to24BitBinary(v int32) string {
	alphabet := [2]byte{'0', '1'}
	out := [24]byte{}
	for i := 23; i >= 0; i-- {
		out[i] = alphabet[v&1]
		v >>= 1
	}
	return bytes.NewBuffer(out[:]).String()
}
