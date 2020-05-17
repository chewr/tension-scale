package state

import (
	"time"

	"github.com/chewr/tension-scale/display"
	"github.com/chewr/tension-scale/display/input"
)

func WaitForInput(required input.Required, received input.Received) display.State {
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

func Work(required input.Required, received input.Received, duration time.Duration) display.State {
	// TODO(rchew) implement
	return nil
}
