package input

import (
	"sync"

	"github.com/chewr/tension-scale/display"
	"github.com/chewr/tension-scale/loadcell"
	"periph.io/x/periph/conn/physic"
)

type EdgeInput interface {
	display.ExpectedInput
	getThresholds() physic.Force
}

type expectedEdgeInputImpl struct {
	minForce physic.Force
}

func (*expectedEdgeInputImpl) GetValue() display.UserInputValue {
	// TODO(rchew) implement usefully
	return nil
}

func (input *expectedEdgeInputImpl) getThresholds() physic.Force {
	return input.minForce
}

var _ display.ActualInput = &DynamicEdgeInput{}

type DynamicEdgeInput struct {
	mu sync.Mutex

	startForce physic.Force
	prevForce  physic.Force
	rising     bool

	maxRisingEdge physic.Force
}

func (input *DynamicEdgeInput) Update(samples ...loadcell.ForceSample) {
	input.mu.Lock()
	defer input.mu.Unlock()
	for _, s := range samples {
		input.rising = s.Force >= input.prevForce
		if !input.rising {
			input.startForce = s.Force
		} else if s.Force-input.startForce > input.maxRisingEdge {
			input.maxRisingEdge = s.Force - input.startForce
		}
		input.prevForce = s.Force
	}
}

func (*DynamicEdgeInput) GetValue() display.UserInputValue {
	// TODO(rchew) implement usefully
	return nil
}

func (input *DynamicEdgeInput) Satisfies(expectedInput display.ExpectedInput) bool {
	input.mu.Lock()
	defer input.mu.Unlock()
	if other, ok := expectedInput.(EdgeInput); ok {
		minForce := other.getThresholds()
		return input.maxRisingEdge >= minForce
	}
	return false
}

func RisingEdge(f physic.Force) EdgeInput {
	return &expectedEdgeInputImpl{
		minForce: f,
	}
}
