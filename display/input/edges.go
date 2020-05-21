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
	getThresholds() (time.Duration, physic.Force)
}

type expectedEdgeInputImpl struct {
	minDuration time.Duration
	minForce    physic.Force
}

func (*expectedEdgeInputImpl) GetVal() display.UserInputValue {
	// TODO(rchew) implement usefully
	return nil
}

func (input *expectedEdgeInputImpl) getThresholds() (time.Duration, physic.Force) {
	return input.minDuration, input.minForce
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

func (*DynamicEdgeInput) GetVal() display.UserInputValue {
	// TODO(rchew) implement usefully
	return nil
}

func (input *DynamicEdgeInput) Satisfies(expectedInput display.ExpectedInput) bool {
	input.mu.Lock()
	defer input.mu.Unlock()
	if other, ok := expectedInput.(EdgeInput); ok {
		// TODO(rchew) move to incremental calculation on Update() to avoid long blocking calls here
		minDuration, minForce := other.getThresholds()
		var (
			startTime  time.Time
			startForce physic.Force
			rising     bool
		)
		for _, s := range input.samples {
			rising = s.Force >= startForce
			if !rising {
				startTime = s.Time
				startForce = s.Force
			}
			if s.Time.Sub(startTime) > minDuration && s.Force-startForce > minForce {
				return true
			}
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
