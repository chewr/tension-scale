package main

import (
	"context"
	"github.com/chewr/tension-scale/led"
	"github.com/chewr/tension-scale/loadcell"
	"github.com/chewr/tension-scale/maxhangs"
	"github.com/chewr/tension-scale/measurement"
	"periph.io/x/periph/conn/gpio"
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
	if err := test(); err != nil {
		panic(err)
	}
	if err := do(); err != nil {
		panic(err)
	}
	if err := test(); err != nil {
		panic(err)
	}
}

func test() error {
	if err := PIN_LED_RED.Out(gpio.High); err != nil {
		return err
	}
	time.Sleep(time.Second / 4)
	if err := PIN_LED_RED.Out(gpio.Low); err != nil {
		return err
	}
	time.Sleep(time.Second / 4)
	if err := PIN_LED_YLW.Out(gpio.High); err != nil {
		return err
	}
	time.Sleep(time.Second / 4)
	if err := PIN_LED_YLW.Out(gpio.Low); err != nil {
		return err
	}
	time.Sleep(time.Second / 4)
	if err := PIN_LED_GRN.Out(gpio.High); err != nil {
		return err
	}
	time.Sleep(time.Second / 4)
	if err := PIN_LED_GRN.Out(gpio.Low); err != nil {
		return err
	}
	return nil
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

	// make tare obvious
	kch := led.Blink(display.YellowOn, display.YellowOff, time.Second/10)
	if err := sensor.Tare(ctx, 10); err != nil {
		return err
	}
	display.YellowOff()
	close(kch)

	input := measurement.NewTimeSeriesDevice(sensor.MeasureTimeSeries(ctx))
	return getWorkout().Run(ctx, display, input)
}

func getWorkout() maxhangs.Workout {
	basicRep := maxhangs.Composite(
		maxhangs.WorkInterval(loadcell.PoundsToRaw(30), 3*time.Second),
		maxhangs.RestInterval(30*time.Second),
		maxhangs.WorkInterval(loadcell.PoundsToRaw(30), 6*time.Second),
		maxhangs.RestInterval(30*time.Second),
		maxhangs.WorkInterval(loadcell.PoundsToRaw(30), 9*time.Second),
	)
	superSetRest := maxhangs.RestInterval(90 * time.Second)

	return maxhangs.Composite(basicRep, superSetRest, basicRep, superSetRest, basicRep)
}
