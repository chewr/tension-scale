package display

import (
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
	baseAbstractState
	isExpiring, isInputDependent bool

	deadline time.Time
	fallback State

	expected ExpectedInput
	actual   ActualInput
}

func (s *stateImpl) InputDependentState() (InputDependentState, bool) {
	if s.isInputDependent {
		return s, true
	}
	return nil, false
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

// TODO(rchew) need some sort of basic metadata for states...
func NewState(opts ...StateBuilderOption) State {
	s := &stateImpl{}
	for _, opt := range opts {
		opt.apply(s)
	}
	return s
}
