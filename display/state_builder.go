package display

import (
	"sync"
	"time"
)

type StateBuilderOption interface {
	apply(builder *stateImpl)
}

type sbOptFn func(builder *stateImpl)

func (fn sbOptFn) apply(builder *stateImpl) {
	fn(builder)
}

func WithExpiryAndFallback(deadline time.Time, fallback State) StateBuilderOption {
	return sbOptFn(func(builder *stateImpl) {
		builder.deadline = deadline
		builder.fallback = fallback
		builder.isExpiring = true
	})
}

func WithExpectedUserInput(expected ExpectedInput, actual ActualInput) StateBuilderOption {
	return sbOptFn(func(builder *stateImpl) {
		builder.expected = expected
		builder.actual = actual
		builder.isInputDependent = true
	})
}

type stateImpl struct {
	// TODO(rchew) finagle a way to get this to be immutable through options
	// ^ not entirely necessary as options are not user-implementable
	stateType WorkoutStateType

	baseAbstractState
	isExpiring, isInputDependent bool

	deadline time.Time
	fallback State

	expected ExpectedInput
	actual   ActualInput

	mutableState *mutableStateImpl
}

type mutableStateImpl struct {
	mu        sync.Mutex
	startTime *time.Time
}

func (s *mutableStateImpl) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.startTime == nil {
		s.startTime = new(time.Time)
		*s.startTime = time.Now()
	}
}

func (s *mutableStateImpl) GetStartTime() time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.startTime != nil {
		return *s.startTime
	}
	// TODO(rchew) leave behavior undefined?
	return time.Time{}
}

func (s *mutableStateImpl) Started() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.startTime != nil
}

func (s *stateImpl) InputDependentState() (InputDependentState, bool) {
	if s.isInputDependent {
		return s, true
	}
	return nil, false
}

func (s *stateImpl) GetType() WorkoutStateType {
	return s.stateType
}

func (s *stateImpl) GetMutableState() MutableState {
	return s.mutableState
}

func (s *stateImpl) InputRequired() ExpectedInput {
	return s.expected
}

func (s *stateImpl) InputReceived() ActualInput {
	return s.actual
}

func (s *stateImpl) Satisfied() bool {
	return s.actual.Satisfies(s.expected)
}

func (s *stateImpl) ExpiringState() (ExpiringState, bool) {
	if s.isExpiring {
		return s, true
	}
	return nil, false
}

func (s *stateImpl) Deadline() time.Time {
	return s.deadline
}

func (s *stateImpl) Fallback() State {
	return s.fallback
}

func NewState(stateType WorkoutStateType, opts ...StateBuilderOption) State {
	s := &stateImpl{
		stateType:    stateType,
		mutableState: &mutableStateImpl{},
	}
	for _, opt := range opts {
		opt.apply(s)
	}
	return s
}
