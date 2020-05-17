package input

import (
	"sync"

	"github.com/chewr/tension-scale/display"
	"periph.io/x/periph/conn/physic"
)

type ForceInput interface {
	display.ExpectedInput
	getForce() physic.Force
}

type instantaneousForceInputImpl struct {
	f physic.Force
}

func (input *instantaneousForceInputImpl) getForce() physic.Force {
	return input.f
}

func (input *instantaneousForceInputImpl) Satisfies(expectedInput display.ExpectedInput) bool {
	if expected, ok := expectedInput.(ForceInput); ok {
		return input.getForce() >= expected.getForce()
	}
	return false
}

var _ display.ActualInput = &DynamicForceInput{}

type DynamicForceInput struct {
	mu sync.Mutex
	f  physic.Force
}

func (input *DynamicForceInput) UpdatingForceInput(f physic.Force) {
	input.mu.Lock()
	defer input.mu.Unlock()
	input.f = f
}

func (input *DynamicForceInput) Satisfies(expected display.ExpectedInput) bool {
	input.mu.Lock()
	defer input.mu.Unlock()
	if f, ok := expected.(ForceInput); ok {
		return input.f >= f.getForce()
	}
	return false
}

func ForceRequired(f physic.Force) display.ExpectedInput {
	return &instantaneousForceInputImpl{f: f}
}

// Deprecated: use DynamicForceInput instead
func ForceReceived(f physic.Force) display.ActualInput {
	return &instantaneousForceInputImpl{f: f}
}
