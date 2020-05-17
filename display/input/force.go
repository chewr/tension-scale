package input

import (
	"github.com/chewr/tension-scale/display"
	"periph.io/x/periph/conn/physic"
)

// TODO(rchew) implement
type Required display.UserInput
type Received display.UserInput

func ForceRequired(f physic.Force) Required {
	// TODO(rchew) implement
	return nil
}

func ForceReceived(f physic.Force) Received {
	// TODO(rchew) implement
	return nil
}

type Edge display.UserInput

// RisingEdge indicates a rising edge force input
func RisingEdge() Edge {
	// TODO(rchew) specify min rising edge size?
	return nil
}
