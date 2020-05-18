package state

import (
	"time"

	"github.com/chewr/tension-scale/display"
)

// TODO(rchew) use opts instead of writing more functions

func WaitForInputWithTimeout(required display.ExpectedInput, received display.ActualInput, deadline time.Time) display.State {
	// TODO(rchew) implement
	return nil
}

func WaitForInput(required display.ExpectedInput, received display.ActualInput) display.State {
	// TODO(rchew) implement
	return nil
}

func Tare(duration time.Duration) display.State {
	// TODO(rchew) implement
	return nil
}

func Rest(duration time.Duration) display.State {
	// TODO(rchew) implement
	return nil
}

func Halt() display.State {
	return nil
}

func Work(required display.ExpectedInput, received display.ActualInput, duration time.Duration) display.State {
	// TODO(rchew) implement
	return nil
}
