package cli

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/chewr/tension-scale/display"
	"github.com/chewr/tension-scale/display/cli/refresh"
	"github.com/chewr/tension-scale/display/input"
	"github.com/fatih/color"
)

type cliDisplay struct {
	mu           sync.Mutex
	currentState display.State
	printer      refresh.Printer
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
			startTime := time.Time{} // TODO(Rchew) acutal time
			// TODO(rchew) error logging
			_ = d.printer.Print(ToCliOutput(startTime, currentState))
		}
	}()
}

func NewCliDisplay(w io.Writer) (display.AutoRefreshingModel, error) {
	// TODO(rchew) wrap
	d := &cliDisplay{
		printer: refresh.NewPrinter(w),
	}
	return d, nil
}

func title(state display.State) refresh.CliOutput {
	// TODO(rchew) move to display pkg
	s := ""
	switch state.GetType() {
	case display.Rest:
		s = "Rest"
	case display.Tare:
		s = "Taring"
	case display.Wait, display.Work:
		s = "Pull"
	case display.Halt:
		fallthrough
	default:
		return refresh.NoShow()
	}
	return refresh.FromString(s)
}

func clock(state display.State) refresh.CliOutput {
	if expiring, ok := state.ExpiringState(); ok {
		ttl := expiring.Deadline().Sub(time.Now())
		return refresh.FromString(fmt.Sprintf("%5.2fs", ttl.Seconds()))
	}
	return refresh.NoShow()
}

// TODO(rchew) should state intrinsically know how long it has been in effect?
func progressBar(startTime time.Time, state display.State) refresh.CliOutput {
	if expiring, ok := state.ExpiringState(); ok {
		deadline := expiring.Deadline()
		totalTime := deadline.Sub(startTime)
		timeElapsed := time.Since(startTime)

		displayColor := color.FgGreen
		if totalTime-timeElapsed < 3*time.Second {
			displayColor = color.FgRed
		}

		// TODO(rchew) flash in middle?

		return refresh.WithColors(
			bar(int64(timeElapsed), int64(totalTime), 20),
			displayColor,
		)
	}
	return refresh.NoShow()
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

func powerBar(state display.State) refresh.CliOutput {
	if dependent, ok := state.InputDependentState(); ok {
		displayColor := color.FgRed
		if dependent.Satisfied() {
			displayColor = color.FgGreen
		}

		// TODO(rchew) this feels clumsy
		forceRequired, requiredOk := dependent.InputRequired().(input.ForceInput)
		forceReceived, receivedOk := dependent.InputReceived().(input.ForceInput)
		if receivedOk && requiredOk {
			if !dependent.Satisfied() &&
				forceReceived.GetForce() > (forceRequired.GetForce()*3)/4 {
				displayColor = color.FgYellow
			}
			s := barWithOverfill(
				int64(forceReceived.GetForce()),
				int64(forceRequired.GetForce()),
				int64(4*forceRequired.GetForce()/3),
				30,
			)
			return refresh.WithColors(s, displayColor)
		}
		// TODO(rchew) print something to reflect whether things have been satisfied?
	}
	return refresh.NoShow()
}

func ToCliOutput(startTime time.Time, state display.State) refresh.CliOutput {
	return refresh.Concat(
		refresh.FromString("    "),
		title(state),
		clock(state),
		progressBar(startTime, state),
		powerBar(state),
	)
}
