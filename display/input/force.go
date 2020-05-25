package input

import (
	"sync"

	"github.com/chewr/tension-scale/display"
	"periph.io/x/periph/conn/physic"
)

type ForceInput interface {
	display.ExpectedInput
	GetForce() physic.Force
}

type instantaneousForceInputImpl struct {
	f physic.Force
}

func (input *instantaneousForceInputImpl) GetForce() physic.Force {
	return input.f
}

func (input *instantaneousForceInputImpl) Satisfies(expectedInput display.ExpectedInput) bool {
	if expected, ok := expectedInput.(ForceInput); ok {
		return input.GetForce() >= expected.GetForce()
	}
	return false
}

func (input *instantaneousForceInputImpl) GetValue() display.UserInputValue {
	return input.f
}

var _ display.ActualInput = &DynamicForceInput{}

type DynamicForceInput struct {
	mu sync.Mutex
	f  physic.Force
}

func (input *DynamicForceInput) UpdateForceInput(f physic.Force) {
	input.mu.Lock()
	defer input.mu.Unlock()
	input.f = f
}

func (input *DynamicForceInput) Satisfies(expected display.ExpectedInput) bool {
	input.mu.Lock()
	defer input.mu.Unlock()
	if f, ok := expected.(ForceInput); ok {
		return input.f >= f.GetForce()
	}
	return false
}

func (input *DynamicForceInput) GetValue() display.UserInputValue {
	return input.f
}

func (input *DynamicForceInput) GetForce() physic.Force {
	input.mu.Lock()
	defer input.mu.Unlock()
	return input.f
}

func ForceRequired(f physic.Force) display.ExpectedInput {
	return &instantaneousForceInputImpl{f: f}
}

// ForceReceived returns an ActualInput for force received
// Deprecated: use DynamicForceInput instead
func ForceReceived(f physic.Force) display.ActualInput {
	return &instantaneousForceInputImpl{f: f}
}
