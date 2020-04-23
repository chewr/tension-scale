package shared

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/chewr/tension-scale/isometric"
	"github.com/chewr/tension-scale/loadcell"
	"github.com/spf13/cobra"
	"periph.io/x/periph/conn/physic"
)

type cliRecorder struct {
	cmd *cobra.Command
}

func CliRecorder(cmd *cobra.Command) isometric.WorkoutRecorder {
	return &cliRecorder{
		cmd: cmd,
	}
}

func (r *cliRecorder) Start(ctx context.Context, descriptor string) (isometric.WorkoutUpdater, error) {
	return &cliWorkoutRecorderUpdater{
		cmd:  r.cmd,
		name: descriptor,
	}, nil
}

type cliWorkoutRecorderUpdater struct {
	mu      sync.Mutex
	cmd     *cobra.Command
	name    string
	samples []loadcell.ForceSample
	closed  bool
}

func (u *cliWorkoutRecorderUpdater) Write(samples ...loadcell.ForceSample) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.closed {
		return isometric.ErrWriteAfterClosed
	}

	u.samples = append(u.samples, samples...)
	return nil
}

func (u *cliWorkoutRecorderUpdater) Finish(outcome isometric.WorkoutOutcome) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.closed {
		return isometric.ErrWriteAfterClosed
	}

	// ensure data is sorted to clean up our output
	sort.Slice(u.samples, func(i, j int) bool {
		return u.samples[i].Time.Before(u.samples[j].Time)
	})

	peakForce := peakForceOverInterval(100*time.Millisecond, u.samples)
	rfd := rateOfForceDevelopment((peakForce*9)/10, u.samples)
	maxForce3s := maxThresholdForceOverInterval(3*time.Second, u.samples)
	maxForce6s := maxThresholdForceOverInterval(6*time.Second, u.samples)
	maxForce9s := maxThresholdForceOverInterval(9*time.Second, u.samples)
	maxForce12s := maxThresholdForceOverInterval(12*time.Second, u.samples)

	sb := new(strings.Builder)
	sb.WriteString(fmt.Sprintf("%s:\n", u.name))
	sb.WriteString(fmt.Sprintf("Peak Force: %s\n", peakForce.String()))
	sb.WriteString(fmt.Sprintf("RFD: %d ms\n", rfd/time.Millisecond))
	if maxForce3s >= physic.Newton {
		sb.WriteString(fmt.Sprintf("Max Force (3s): %s\n", maxForce3s.String()))
	}
	if maxForce6s >= physic.Newton {
		sb.WriteString(fmt.Sprintf("Max Force (6s): %s\n", maxForce6s.String()))
	}
	if maxForce9s >= physic.Newton {
		sb.WriteString(fmt.Sprintf("Max Force (9s): %s\n", maxForce9s.String()))
	}
	if maxForce12s >= physic.Newton {
		sb.WriteString(fmt.Sprintf("Max Force (12s): %s\n", maxForce12s.String()))
	}

	u.closed = true
	_, err := fmt.Fprintln(u.cmd.OutOrStdout(), sb.String())
	return err
}

func peakForceOverInterval(d time.Duration, samples []loadcell.ForceSample) physic.Force {
	weight := func(i int) int64 {
		// adjust weights
		if i > 0 {
			// weights are in time units rounded to microseconds to avoid int64 overflow later
			return int64(samples[i].Time.Sub(samples[i-1].Time) / time.Microsecond)
		}
		return 0
	}
	movWgtTot := int64(0)
	movWgtSum := int64(0)
	maxMovWgtAvg := physic.Force(0)
	startPtr := 0
	for i, s := range samples {
		movWgtSum += weight(i) * int64(s.Force)
		movWgtTot += weight(i)

		// set startPtr
		for ; s.Sub(samples[startPtr].Time) >= d; startPtr++ {
			movWgtSum -= int64(samples[startPtr].Force) * weight(startPtr)
			movWgtTot -= weight(startPtr)
		}

		// wait until we've processed enough data to calculate moving averages
		if startPtr > 0 {
			movWgtAvg := physic.Force(movWgtSum / movWgtTot)
			if movWgtAvg > maxMovWgtAvg {
				maxMovWgtAvg = movWgtAvg
			}
		}
	}
	return maxMovWgtAvg
}

func maxThresholdForceOverInterval(d time.Duration, samples []loadcell.ForceSample) physic.Force {
	startPtr := 0
	thresholdForce := physic.Force(0)
	for i, s := range samples {
		// set startPtr
		for ; s.Sub(samples[startPtr].Time) >= d; startPtr++ {
		}

		// wait until we've processed enough data to fill a sliding window
		if startPtr > 0 {
			// O(n^2) for now; sucks to suck
			minForceInWindow := s.Force
			for j := startPtr; j <= i; j++ {
				if samples[j].Force < minForceInWindow {
					minForceInWindow = samples[j].Force
				}
			}
			if minForceInWindow > thresholdForce {
				thresholdForce = minForceInWindow
			}
		}
	}
	return thresholdForce
}

func rateOfForceDevelopment(f physic.Force, samples []loadcell.ForceSample) time.Duration {
	smoothed := samples
	dFdT := make([]int64, len(smoothed))
	dFdT[0] = 0
	for i := 1; i < len(smoothed)-1; i++ {
		dFdT[i] = int64(smoothed[i+1].Force-smoothed[i-1].Force) / int64(smoothed[i+1].Time.Sub(smoothed[i-1].Time))
	}
	dFdT[len(smoothed)-1] = 0

	curRisingLen := 0
	maxRisingLen := 0
	maxRiseStart := 0
	for i, d := range dFdT {
		if d > 0 {
			curRisingLen++
		}
		if d <= 0 {
			if curRisingLen > maxRisingLen {
				maxRisingLen = curRisingLen
				maxRiseStart = i - curRisingLen
			}
			curRisingLen = 0
		}
	}

	maxForce := physic.Force(0)
	timeToMaxForce := time.Duration(0)
	for i := maxRiseStart; i < len(samples); i++ {
		if samples[i].Force >= f {
			return samples[i].Time.Sub(samples[maxRiseStart].Time)
		}
		if samples[i].Force > maxForce {
			maxForce = samples[i].Force
			timeToMaxForce = samples[i].Time.Sub(samples[maxRiseStart].Time)
		}
	}

	// if for some reason threshold f was never reached,
	// just return time to Max Force so we have something
	// to report
	return timeToMaxForce
}

func (u *cliWorkoutRecorderUpdater) Close() {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.closed = true
}
