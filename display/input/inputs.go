package input

import (
	"github.com/chewr/tension-scale/display"
)

type noInput struct{}

func (noInput) Satisfies(_ display.ExpectedInput) bool { return false }
func (noInput) GetVal() display.UserInputValue         { return nil }

func None() display.ActualInput {
	return noInput{}
}
