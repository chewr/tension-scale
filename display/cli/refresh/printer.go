package refresh

import (
	"encoding/hex"
	"fmt"
	"io"
	"sync"
	"unicode/utf8"
)

type CliOutput interface {
	NoColor() string
	WithColor() string
}

type Printer interface {
	Print(output CliOutput) error
}

func NewPrinter(w io.Writer) Printer {
	return &refreshingWriter{
		writer: w,
	}
}

type refreshingWriter struct {
	mu         sync.Mutex
	writer     io.Writer
	lastOutput string
}

func (w *refreshingWriter) erase() error {
	n := utf8.RuneCountInString(w.lastOutput)
	del, _ := hex.DecodeString("7f")
	for _, c := range []string{"\b", string(del)} {
		for i := 0; i < n; i++ {
			if _, err := fmt.Fprintf(w.writer, c); err != nil {
				return err
			}
		}
	}
	// erases to end of line
	if _, err := fmt.Fprintf(w.writer, "\r\033[K"); err != nil {
		return err
	}
	w.lastOutput = ""
	return nil
}

func (s *refreshingWriter) Print(output CliOutput) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.erase(); err != nil {
		return err
	}

	outColor := fmt.Sprintf("%s ", output.WithColor())
	if _, err := fmt.Fprint(s.writer, outColor); err != nil {
		return err
	}
	s.lastOutput = fmt.Sprintf("%s ", output.NoColor())
	return nil
}
