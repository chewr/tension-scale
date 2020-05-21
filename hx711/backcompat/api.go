package backcompat

import (
	"time"

	"periph.io/x/periph/experimental/conn/analog"
	periphimpl "periph.io/x/periph/experimental/devices/hx711"
)

// HX711 replicates the interface of the periph.io driver for the
// HX711 load cell amplifier. This interface is deprecated, use V2 instead
//
// Deprecated: use V2 instead
type HX711 interface {
	analog.PinADC
	ReadContinuous() <-chan analog.Sample
	ReadTimeout(timeout time.Duration) (int32, error)
	SetInputMode(inputMode periphimpl.InputMode) error
	IsReady() bool
}
