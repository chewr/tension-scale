package main

import (
	"fmt"
	"os"
	"time"

	"github.com/chewr/tension-scale/errutil"
	"periph.io/x/periph/experimental/devices/hx711"
	"periph.io/x/periph/host"
	"periph.io/x/periph/host/rpi"
)

func main() {
	if err := mainE(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func mainE() error {
	if _, err := host.Init(); err != nil {
		return err
	}

	// HX711 board configured with CLK @ P1_31 and DOUT @ P1_29
	hx, err := hx711.New(rpi.P1_31, rpi.P1_29)
	if err != nil {
		return err
	}

	time.AfterFunc(10*time.Second, func() { errutil.PanicOnErr(hx.Halt()) })

	sampleCh := hx.ReadContinuous()
	for s := range sampleCh {
		// observe that many samples are -1
		fmt.Printf("%8d\n", s.Raw)
	}
	return nil
}
