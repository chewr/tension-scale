package state

import (
	"time"

	"github.com/chewr/tension-scale/display"
)

// TODO(rchew) use opts instead of writing more functions

func WaitForInputWithTimeout(required display.ExpectedInput, received display.ActualInput, deadline time.Time) display.State {
	return display.NewState(
		display.WithExpectedUserInput(required, received),
		display.WithExpiryAndFallback(deadline, Halt()),
	)
}

func WaitForInput(required display.ExpectedInput, received display.ActualInput) display.State {
	return display.NewState(display.WithExpectedUserInput(required, received))
}

func Tare(duration time.Duration) display.State {
	// TODO(rchew) take deadline instead of ttl
	deadline := time.Now().Add(duration)
	// TODO(rchew) better to fall back to original state rather than halt?
	return display.NewState(display.WithExpiryAndFallback(deadline, Halt()))
}

func Rest(duration time.Duration) display.State {
	// TODO(rchew) take deadline instead of ttl
	deadline := time.Now().Add(duration)
	// TODO(rchew) better to fall back to original state rather than halt?
	return display.NewState(display.WithExpiryAndFallback(deadline, Halt()))
}

func Halt() display.State {
	return display.NewState()
}

func Work(required display.ExpectedInput, received display.ActualInput, duration time.Duration) display.State {
	// TODO(rchew) take deadline instead of ttl
	deadline := time.Now().Add(duration)
	return display.NewState(
		display.WithExpectedUserInput(required, received),
		// TODO(rchew) why is Halt always the fallback? What should the default state be?
		display.WithExpiryAndFallback(deadline, Halt()),
	)
}
