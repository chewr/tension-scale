package cli

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/chewr/tension-scale/display"
	"github.com/gosuri/uilive"
)

type cliDisplay struct {
	mu           sync.Mutex
	currentState display.State
	uw           *uilive.Writer
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

func toCliOutput(state display.State) string {
	// TODO(rchew) implement
	return ""
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
			_, _ = fmt.Fprintln(d.uw, toCliOutput(currentState))
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
