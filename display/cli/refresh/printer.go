package refresh

import (
	"encoding/hex"
	"fmt"
	"github.com/fatih/color"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
	"unicode/utf8"
)

type CliOutput interface {
	NoColor() string
	WithColor() string
}

// shamelessly stolen from briandowns/spinner
// Spinner struct to hold the provided options.
type Spinner struct {
	mu         *sync.RWMutex //
	Delay      time.Duration // Delay is the speed of the indicator
	output     CliOutput
	chars      []string                      // chars holds the chosen character set
	FinalMSG   string                        // string displayed after Stop() is called
	lastOutput string                        // last character(set) written
	Writer     io.Writer                     // to make testing better, exported so users have access. Use `WithWriter` to update after initialization.
	active     bool                          // active holds the state of the spinner
	stopChan   chan struct{}                 // stopChan is a channel used to stop the indicator
	HideCursor bool                          // hideCursor determines if the cursor is visible
	PreUpdate  func(s *Spinner)              // will be triggered before every spinner update
	PostUpdate func(s *Spinner)              // will be triggered after every spinner update
}

// New provides a pointer to an instance of Spinner with the supplied options.
func New(output CliOutput, d time.Duration, options ...Option) *Spinner {
	s := &Spinner{
		Delay:    d,
		output:output,
		mu:       &sync.RWMutex{},
		Writer:   color.Output,
		active:   false,
		stopChan: make(chan struct{}, 1),
	}

	for _, option := range options {
		option(s)
	}
	return s
}

// Option is a function that takes a spinner and applies
// a given configuration.
type Option func(*Spinner)

// Options contains fields to configure the spinner.
type Options struct {
	FinalMSG   string
	HideCursor bool
}

// WithFinalMSG adds the given string ot the spinner
// as the final message to be written.
func WithFinalMSG(finalMsg string) Option {
	return func(s *Spinner) {
		s.FinalMSG = finalMsg
	}
}

// WithHiddenCursor hides the cursor
// if hideCursor = true given.
func WithHiddenCursor(hideCursor bool) Option {
	return func(s *Spinner) {
		s.HideCursor = hideCursor
	}
}

// WithWriter adds the given writer to the spinner. This
// function should be favored over directly assigning to
// the struct value.
func WithWriter(w io.Writer) Option {
	return func(s *Spinner) {
		s.mu.Lock()
		s.Writer = w
		s.mu.Unlock()
	}
}

// Active will return whether or not the spinner is currently active.
func (s *Spinner) Active() bool {
	return s.active
}

// Start will start the indicator.
func (s *Spinner) Start() {
	s.mu.Lock()
	if s.active {
		s.mu.Unlock()
		return
	}
	if s.HideCursor && runtime.GOOS != "windows" {
		// hides the cursor
		fmt.Print("\033[?25l")
	}
	s.active = true
	s.mu.Unlock()

	go func() {
		for {
			select {
			case <-s.stopChan:
				return
			default:
				if !s.active {
					return
				}
				s.mu.Lock()
				s.erase()

				if s.PreUpdate != nil {
					s.PreUpdate(s)
				}

				var outColor string
				if runtime.GOOS == "windows" {
					if s.Writer == os.Stderr {
						outColor = fmt.Sprintf("\r%s ", s.output.NoColor())
					} else {
						outColor = fmt.Sprintf("\r%s ", s.output.WithColor())
					}
				} else {
					outColor = fmt.Sprintf("%s ", s.output.WithColor())
				}
				outPlain := fmt.Sprintf("%s ", s.output.NoColor())
				fmt.Fprint(s.Writer, outColor)
				s.lastOutput = outPlain
				delay := s.Delay

				if s.PostUpdate != nil {
					s.PostUpdate(s)
				}

				s.mu.Unlock()
				time.Sleep(delay)
			}
		}
	}()
}

// Stop stops the indicator.
func (s *Spinner) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.active {
		s.active = false
		if s.HideCursor && runtime.GOOS != "windows" {
			// makes the cursor visible
			fmt.Print("\033[?25h")
		}
		s.erase()
		if s.FinalMSG != "" {
			fmt.Fprintf(s.Writer, s.FinalMSG)
		}
		s.stopChan <- struct{}{}
	}
}

// Restart will stop and start the indicator.
func (s *Spinner) Restart() {
	s.Stop()
	s.Start()
}

// UpdateSpeed will set the indicator delay to the given value.
func (s *Spinner) UpdateSpeed(d time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Delay = d
}

// erase deletes written characters.
// Caller must already hold s.lock.
func (s *Spinner) erase() {
	n := utf8.RuneCountInString(s.lastOutput)
	if runtime.GOOS == "windows" {
		clearString := "\r"
		for i := 0; i < n; i++ {
			clearString += " "
		}
		clearString += "\r"
		fmt.Fprintf(s.Writer, clearString)
		s.lastOutput = ""
		return
	}
	del, _ := hex.DecodeString("7f")
	for _, c := range []string{"\b", string(del)} {
		for i := 0; i < n; i++ {
			fmt.Fprintf(s.Writer, c)
		}
	}
	fmt.Fprintf(s.Writer, "\r\033[K") // erases to end of line
	s.lastOutput = ""
}

// Lock allows for manual control to lock the spinner.
func (s *Spinner) Lock() {
	s.mu.Lock()
}

// Unlock allows for manual control to unlock the spinner.
func (s *Spinner) Unlock() {
	s.mu.Unlock()
}
