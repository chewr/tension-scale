package display

import (
	"context"
	"time"
)

type AbstractState interface{}

type State interface {
	AbstractState
	ExpiringState() (ExpiringState, bool)
	InputDependentState() (InputDependentState, bool)
}

type InputDependentState interface {
	AbstractState
	InputRequired() UserInput
	InputReceived() UserInput
	Satisfied() bool
}

type ExpiringState interface {
	AbstractState
	Deadline() time.Time
	Fallback() State
}

type UserInput interface{}

type Model interface {
	UpdateState(state State) error
}

type RunningModel interface {
	Model
	Run(ctx context.Context)
}
