package main

import (
	"fmt"
	"os"
	"periph.io/x/periph/experimental/devices/hx711"
	"periph.io/x/periph/host"
	"periph.io/x/periph/host/rpi"
	"time"
)

func main() {
	host.Init()
	if err := mainE(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func mainE() error {
	// HX711 board configured with CLK @ P1_31 and DOUT @ P1_29
	hx, err := hx711.New(rpi.P1_31, rpi.P1_29)
	if err != nil {
		return err
	}

	time.AfterFunc(10*time.Second, func() { hx.Halt() })

	sampleCh := hx.ReadContinuous()
	for s := range sampleCh {
		// observe that many samples are -1
		fmt.Println("%8d", s.Raw)
	}
	return nil
}
