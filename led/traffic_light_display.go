package led

import (
	"errors"

	"github.com/chewr/tension-scale/display"
)

func (*TrafficLight) UpdateState(state display.State) error {
	return errors.New("Not implemented")
}
