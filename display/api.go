package display

import (
	"context"
	"fmt"
	"time"
)

type WorkoutStateType int

const (
	Halt = iota
	Rest
	Work
	Tare
	Wait
)

type AbstractState interface {
	// TODO(rchew) decide if this is necessary
	noImplementAbstractState()
	GetType() WorkoutStateType
}

type baseAbstractState struct{}

func (baseAbstractState) noImplementAbstractState() {}

type State interface {
	AbstractState
	// TODO(rchew) better to just do this with casting?
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
	GetVal() UserInputValue
}

type UserInputValue interface {
	fmt.Stringer
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
