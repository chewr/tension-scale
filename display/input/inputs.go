package input

import (
	"github.com/chewr/tension-scale/display"
)

type functionalUserInput func(display.ExpectedInput) bool

func (fn functionalUserInput) Satisfies(input display.ExpectedInput) bool {
	return fn(input)
}

func None() display.ActualInput {
	return functionalUserInput(func(_ display.ExpectedInput) bool {
		return false
	})
}
