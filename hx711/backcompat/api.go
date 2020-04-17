package backcompat

import (
	"periph.io/x/periph/experimental/conn/analog"
	periphimpl "periph.io/x/periph/experimental/devices/hx711"
	"time"
)

// deprecated, use V2 instead
type HX711 interface {
	analog.PinADC
	ReadContinuous() <-chan analog.Sample
	ReadTimeout(timeout time.Duration) (int32, error)
	SetInputMode(inputMode periphimpl.InputMode) error
	IsReady() bool
}
