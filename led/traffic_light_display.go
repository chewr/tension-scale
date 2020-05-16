package led

import (
	"errors"

	"github.com/chewr/tension-scale/display"
)

func (*TrafficLight) UpdateState(state display.State) error {
	if state == nil {
		return errors.New("Unrecognized state")
	}
	return errors.New("Not implemented")
}
