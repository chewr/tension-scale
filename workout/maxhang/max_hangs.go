package maxhang

import (
	"errors"
	"time"

	"github.com/chewr/tension-scale/isometric"
	"periph.io/x/periph/conn/physic"
)

type Week int

const (
	Week1 Week = 1
	Week2 Week = 2
	Week3 Week = 3
	Week4 Week = 4
)

var ErrWeekOutOfRange = errors.New("week out of range: max hangs are a four week cycle")

func MaxHangWorkout(week Week, weight physic.Force) (isometric.Workout, error) {
	var (
		rep  isometric.Workout
		sets int
	)
	switch week {
	default:
		return nil, ErrWeekOutOfRange
	case Week3:
		rep = isometric.Composite(
			isometric.WorkInterval(weight, 3*time.Second),
			isometric.RestInterval(30*time.Second),
			isometric.WorkInterval(weight, 6*time.Second),
			isometric.RestInterval(30*time.Second),
			isometric.WorkInterval(weight, 9*time.Second),
		)
		sets = 5
	case Week2:
		rep = isometric.Composite(
			isometric.WorkInterval(weight, 3*time.Second),
			isometric.RestInterval(30*time.Second),
			isometric.WorkInterval(weight, 6*time.Second),
			isometric.RestInterval(30*time.Second),
			isometric.WorkInterval(weight, 9*time.Second),
		)
		sets = 4
	case Week1:
		rep = isometric.Composite(
			isometric.WorkInterval(weight, 3*time.Second),
			isometric.RestInterval(30*time.Second),
			isometric.WorkInterval(weight, 6*time.Second),
			isometric.RestInterval(30*time.Second),
			isometric.WorkInterval(weight, 9*time.Second),
		)
		sets = 3
	case Week4:
		rep = isometric.Composite(
			isometric.WorkInterval(weight, 3*time.Second),
			isometric.RestInterval(30*time.Second),
			isometric.WorkInterval(weight, 6*time.Second),
			isometric.RestInterval(30*time.Second),
			isometric.WorkInterval(weight, 9*time.Second),
			isometric.RestInterval(30*time.Second),
			isometric.WorkInterval(weight, 12*time.Second),
		)
		sets = 3
	}
	var supersets []isometric.Workout
	supersets = append(supersets, isometric.SetupInterval(time.Minute), rep)
	for i := 0; i < sets; i++ {
		supersets = append(supersets, isometric.RestInterval(90*time.Second), rep)
	}
	return isometric.Composite(supersets...), nil
}
