package main

import (
	"context"
	"github.com/chewr/tension-scale/led"
	"github.com/chewr/tension-scale/loadcell"
	"github.com/chewr/tension-scale/maxhangs"
	"github.com/chewr/tension-scale/measurement"
	"periph.io/x/periph/experimental/devices/hx711"
	"periph.io/x/periph/host/rpi"
	"time"
)

var (
	PIN_LED_RED = rpi.P1_21 // GPIO9
	PIN_LED_YLW = rpi.P1_19 // GPIO10
	PIN_LED_GRN = rpi.P1_23 // GPIO11

	PIN_HX711_SCLK = rpi.P1_31 // GPIO6
	PIN_HX711_DOUT = rpi.P1_29 // GPIO5
)

func main() {
	if err := do(); err != nil {
		panic(err)
	}
}

func do() error {
	dev, err := hx711.New(PIN_HX711_SCLK, PIN_HX711_DOUT)
	if err != nil {
		return err
	}
	sensor := loadcell.New(dev)
	display := led.TrafficLight{
		Green:  PIN_LED_GRN,
		Yellow: PIN_LED_YLW,
		Red:    PIN_LED_RED,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := sensor.Tare(ctx, 0); err != nil {
		return err
	}

	input := measurement.NewTimeSeriesDevice(sensor.MeasureTimeSeries(ctx))
	return getWorkout().Run(ctx, display, input)
}

func getWorkout() maxhangs.Workout {
	basicRep := maxhangs.Composite(
		maxhangs.WorkInterval(loadcell.PoundsToRaw(80), 3*time.Second),
		maxhangs.RestInterval(30*time.Second),
		maxhangs.WorkInterval(loadcell.PoundsToRaw(80), 6*time.Second),
		maxhangs.RestInterval(30*time.Second),
		maxhangs.WorkInterval(loadcell.PoundsToRaw(80), 9*time.Second),
	)
	superSetRest := maxhangs.RestInterval(90 * time.Second)

	return maxhangs.Composite(basicRep, superSetRest, basicRep, superSetRest, basicRep)
}
