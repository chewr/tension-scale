package display

import (
	"context"
	"time"
)

type AbstractState interface {
	// TODO(rchew) decide if this is necessary
	noImplementAbstractState()
}

type baseAbstractState struct{}

func (baseAbstractState) noImplementAbstractState() {}

type State interface {
	AbstractState
	ExpiringState() (ExpiringState, bool)
	InputDependentState() (InputDependentState, bool)
}

type InputDependentState interface {
	AbstractState
	InputRequired() ExpectedInput
	InputReceived() ActualInput
	Satisfied() bool
}

type ExpiringState interface {
	AbstractState
	Deadline() time.Time
	Fallback() State
}

type UserInput interface {
}

type ExpectedInput interface {
	UserInput
}

type ActualInput interface {
	UserInput
	Satisfies(input ExpectedInput) bool
}

type Model interface {
	UpdateState(state State) error
}

type AutoRefreshingModel interface {
	Model
	Start(ctx context.Context)
}
