package backcompat

import (
	"time"

	"periph.io/x/periph/experimental/conn/analog"
	periphimpl "periph.io/x/periph/experimental/devices/hx711"
)

// deprecated, use V2 instead
// HX711 replicates the interface of the periph.io driver for the
// HX711 load cell amplifier
type HX711 interface {
	analog.PinADC
	ReadContinuous() <-chan analog.Sample
	ReadTimeout(timeout time.Duration) (int32, error)
	SetInputMode(inputMode periphimpl.InputMode) error
	IsReady() bool
}
