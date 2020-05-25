package cli

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/chewr/tension-scale/display"
	"github.com/chewr/tension-scale/display/input"
	"github.com/fatih/color"
	"github.com/gosuri/uilive"
)

type cliDisplay struct {
	mu           sync.Mutex
	currentState display.State
	// TODO(rchew) uilive is trash
	uw *uilive.Writer
}

func (d *cliDisplay) UpdateState(state display.State) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.currentState = state
	return nil
}

func (d *cliDisplay) getCurrentState() display.State {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.currentState == nil {
		return nil
	}

	if expiring, ok := d.currentState.ExpiringState(); ok {
		if time.Now().After(expiring.Deadline()) {
			d.currentState = expiring.Fallback()
		}
	}

	return d.currentState
}

func (d *cliDisplay) Start(ctx context.Context) {
	t := time.NewTicker(50 * time.Millisecond)
	go func() {
		defer t.Stop()
		for range t.C {
			select {
			case <-ctx.Done():
				return
			default:
			}
			currentState := d.getCurrentState()
			// TODO(rchew) error logging
			startTime := time.Time{} // TODO(Rchew) acutal time
			_, _ = fmt.Fprintln(d.uw, ToCliOutput(startTime, currentState))
			_ = d.uw.Flush()
		}
	}()
}

func NewCliDisplay(w io.Writer) (display.AutoRefreshingModel, error) {
	// TODO(rchew) wrap
	uw := uilive.New()
	uw.Out = w
	d := &cliDisplay{
		uw: uw,
	}
	return d, nil
}

////
// TODO(rchew) Clean up this implementation:
// - Apply states to build a cli output type
//   - cli output type has colored and uncolored output

// TODO(rchew) multi-component ui

func title(state display.State) string {
	// TODO(rchew) move to display pkg
	switch state.GetType() {
	case display.Rest:
		return "Rest"
	case display.Tare:
		return "Taring"
	case display.Wait, display.Work:
		return "Pull"
	case display.Halt:
		fallthrough
	default:
		return ""
	}
}

func clock(state display.State) string {
	if expiring, ok := state.ExpiringState(); ok {
		ttl := expiring.Deadline().Sub(time.Now())
		return fmt.Sprintf("%5.2fs", ttl.Seconds())
	}
	return ""
}

// TODO(rchew) should state intrinsically know how long it has been in effect?
func progressBar(startTime time.Time, state display.State) string {
	if expiring, ok := state.ExpiringState(); ok {
		deadline := expiring.Deadline()
		totalTime := deadline.Sub(startTime)
		timeElapsed := time.Since(startTime)

		return bar(int64(timeElapsed), int64(totalTime), 20)
	}
	return ""
}

func bar(val, max int64, width int) string {
	filledBars := int(float64(width) * float64(val) / float64(max))
	b := strings.Builder{}
	for i := 0; i < width; i++ {
		if i < filledBars {
			b.WriteRune('█')
		} else {
			b.WriteRune('░')
		}
	}
	return b.String()
}

func barWithOverfill(val, threshold, overfill int64, width int) string {
	filledBars := int(int64(width-1) * val / overfill)
	thresholdBar := int(int64(width-1) * threshold / overfill)
	b := strings.Builder{}
	for i := 0; i < width; i++ {
		if i < filledBars {
			b.WriteRune('█')
		} else {
			b.WriteRune('░')
		}
		if i == thresholdBar-1 {
			b.WriteRune('|')
		}
	}
	return b.String()
}

func powerBar(state display.State) string {
	if dependent, ok := state.InputDependentState(); ok {
		displayColor := color.New(color.FgRed)
		if dependent.Satisfied() {
			displayColor = color.New(color.FgGreen)
		}

		// TODO(rchew) this feels clumsy
		forceRequired, requiredOk := dependent.InputRequired().(input.ForceInput)
		forceReceived, receivedOk := dependent.InputReceived().(input.ForceInput)
		if receivedOk && requiredOk {
			if !dependent.Satisfied() &&
				forceReceived.GetForce() > (forceRequired.GetForce()*3)/4 {
				displayColor = color.New(color.FgYellow)
			}
			return displayColor.Sprint(barWithOverfill(
				int64(forceReceived.GetForce()),
				int64(forceRequired.GetForce()),
				int64(4*forceRequired.GetForce()/3),
				30,
			))
		}
		return fmt.Sprintf("Required force: %v, Received Force: %v", requiredOk, receivedOk)
	}
	return ""
}

func ToCliOutput(startTime time.Time, state display.State) string {
	return strings.Join([]string{
		title(state),
		clock(state),
		progressBar(startTime, state),
		powerBar(state),
	}, "    ")
}
