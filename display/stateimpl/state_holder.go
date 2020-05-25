package stateimpl

import (
	"errors"
	"sync"
	"time"

	"github.com/chewr/tension-scale/display"
)

var ErrStateUninitialized = errors.New("no state currently set")

// TODO(rchew) names...
type StateHolder struct {
	mu           sync.Mutex
	currentState display.State
}

func (h *StateHolder) setState(state display.State) {
	h.currentState = state
	h.currentState.GetMutableState().Start()
}

func (h *StateHolder) UpdateState(state display.State) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.setState(state)
	return nil
}

func (h *StateHolder) GetCurrentState() (display.State, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.currentState == nil {
		return nil, ErrStateUninitialized
	}

	if state, expiring := h.currentState.ExpiringState(); expiring {
		if time.Now().After(state.Deadline()) {
			h.setState(state.Fallback())
		}
	}

	return h.currentState, nil
}
