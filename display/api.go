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

func (t WorkoutStateType) String() string {
	switch t {
	case Halt:
		return "Halt"
	case Rest:
		return "Rest"
	case Work:
		return "Pull"
	case Tare:
		return "Taring"
	case Wait:
		return "Ready"
	default:
		return "Unknown"
	}
}

type AbstractState interface {
	// TODO(rchew) decide if this is necessary
	noImplementAbstractState()
	GetType() WorkoutStateType
	GetMutableState() MutableState
}

type MutableState interface {
	Start()
	Started() bool
	GetStartTime() time.Time
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
	// TODO(rchew): make Satisfied return read-once channel and add an evolving state
	Satisfied() bool
}

type ExpiringState interface {
	AbstractState
	Deadline() time.Time
	Fallback() State
}

type UserInput interface {
	GetValue() UserInputValue
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
