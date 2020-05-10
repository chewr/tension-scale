package state

import "github.com/chewr/tension-scale/display"

var (
	Waiting                 display.State = nil
	Taring                  display.State = nil
	Resting                 display.State = nil
	PullingHardEnough       display.State = nil
	PullingButNotHardEnough display.State = nil
)
