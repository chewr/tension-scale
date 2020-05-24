package cli

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/chewr/tension-scale/display"
)

type cliDisplay struct {
	s            *spinner.Spinner
	mu           sync.Mutex
	currentState display.State
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
	d.s.Start()
	go func() {
		<-ctx.Done()
		d.s.Stop()
	}()

	go func() {
		for {
			currentState := d.getCurrentState()
			d.s.UpdateCharSet()
			d.s.Color()
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()
}

func NewCliDisplay(w io.Writer) (display.AutoRefreshingModel, error) {
	var (
		// among the default charsets, [14] is among the few which look good at 50ms refresh rate
		defaultCharset  = []string{"waiting..."}
		refreshInterval = 50 * time.Millisecond
	)
	s := spinner.New(
		defaultCharset,
		refreshInterval,
		spinner.WithWriter(w),
		spinner.WithHiddenCursor(true),
	)
	d := &cliDisplay{
		s: s,
	}
	return d, nil
}
