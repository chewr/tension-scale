package input

import (
	"github.com/chewr/tension-scale/display"
)

type noInput struct{}

func (noInput) Satisfies(_ display.ExpectedInput) bool { return false }
func (noInput) GetValue() display.UserInputValue       { return nil }

// None returns the zero value of input; it never satisfies anything
//
// Deprecated: use an actual input
func None() display.ActualInput {
	return noInput{}
}
