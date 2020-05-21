package input

import (
	"sync"
	"time"

	"github.com/chewr/tension-scale/display"
	"github.com/chewr/tension-scale/loadcell"
	"periph.io/x/periph/conn/physic"
)

type EdgeInput interface {
	display.ExpectedInput
	getThresholds() physic.Force
}

type expectedEdgeInputImpl struct {
	minDuration time.Duration
	minForce    physic.Force
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
	mu      sync.Mutex
	samples []loadcell.ForceSample
}

func (input *DynamicEdgeInput) Update(samples ...loadcell.ForceSample) {
	input.mu.Lock()
	defer input.mu.Unlock()
	input.samples = append(input.samples, samples...)
}

func (*DynamicEdgeInput) GetValue() display.UserInputValue {
	// TODO(rchew) implement usefully
	return nil
}

func (input *DynamicEdgeInput) Satisfies(expectedInput display.ExpectedInput) bool {
	input.mu.Lock()
	defer input.mu.Unlock()
	if other, ok := expectedInput.(EdgeInput); ok {
		// TODO(rchew) move to incremental calculation on Update() to avoid long blocking calls here
		minForce := other.getThresholds()
		var (
			startForce physic.Force
			prevForce  physic.Force
			rising     bool
		)
		for _, s := range input.samples {
			rising = s.Force >= prevForce
			if !rising {
				startForce = s.Force
			}
			if s.Force-startForce > minForce {
				return true
			}
			prevForce = s.Force
		}
	}
	return false
}

func RisingEdge(t time.Duration, f physic.Force) EdgeInput {
	return &expectedEdgeInputImpl{
		minDuration: t,
		minForce:    f,
	}
}
