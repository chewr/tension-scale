package maxhang

import (
	"errors"
	"time"

	"github.com/chewr/tension-scale/isometric"
	"github.com/chewr/tension-scale/isometric/interval"
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
		rep = interval.Composite(
			interval.WorkInterval(weight, 3*time.Second),
			interval.RestInterval(30*time.Second),
			interval.WorkInterval(weight, 6*time.Second),
			interval.RestInterval(30*time.Second),
			interval.WorkInterval(weight, 9*time.Second),
		)
		sets = 5
	case Week2:
		rep = interval.Composite(
			interval.WorkInterval(weight, 3*time.Second),
			interval.RestInterval(30*time.Second),
			interval.WorkInterval(weight, 6*time.Second),
			interval.RestInterval(30*time.Second),
			interval.WorkInterval(weight, 9*time.Second),
		)
		sets = 4
	case Week1:
		rep = interval.Composite(
			interval.WorkInterval(weight, 3*time.Second),
			interval.RestInterval(30*time.Second),
			interval.WorkInterval(weight, 6*time.Second),
			interval.RestInterval(30*time.Second),
			interval.WorkInterval(weight, 9*time.Second),
		)
		sets = 3
	case Week4:
		rep = interval.Composite(
			interval.WorkInterval(weight, 3*time.Second),
			interval.RestInterval(30*time.Second),
			interval.WorkInterval(weight, 6*time.Second),
			interval.RestInterval(30*time.Second),
			interval.WorkInterval(weight, 9*time.Second),
			interval.RestInterval(30*time.Second),
			interval.WorkInterval(weight, 12*time.Second),
		)
		sets = 3
	}
	var supersets []isometric.Workout
	supersets = append(supersets, interval.SetupInterval(time.Minute), rep)
	for i := 1; i < sets; i++ {
		supersets = append(supersets, interval.RestInterval(90*time.Second), rep)
	}
	return interval.Composite(supersets...), nil
}
