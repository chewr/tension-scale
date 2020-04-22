package shared

import (
	"github.com/chewr/tension-scale/hx711"
	"github.com/chewr/tension-scale/isometric"
	"github.com/chewr/tension-scale/led"
	"github.com/chewr/tension-scale/loadcell"
	"periph.io/x/periph/host/rpi"
)

// TODO(rchew): make reusable
// TODO(rchew): make portable
// TODO(rchew): make configurable

func SetupDisplay() (*led.TrafficLight, error) {
	grn := rpi.P1_23
	ylw := rpi.P1_19
	red := rpi.P1_21
	return led.NewTrafficLight(grn, ylw, red)
}

func SetupLoadCell() (loadcell.Sensor, error) {
	hx, err := hx711.New(rpi.P1_31, rpi.P1_29)
	if err != nil {
		return nil, err
	}
	return loadcell.NewHx711(hx, loadcell.TrueSun400Slow), nil
}

func SetupOutput() (isometric.WorkoutRecorder, error) {
	const defaultOutputDir = "~/Documents/Workouts"
	return isometric.NewCsvFileRecorder(defaultOutputDir)
}
