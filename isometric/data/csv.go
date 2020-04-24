package data

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/chewr/tension-scale/isometric"
	"github.com/chewr/tension-scale/loadcell"
	"periph.io/x/periph/conn/physic"
)

type csvFileRecorder struct {
	dir string
}

type csvFileRecorderUpdater struct {
	mu       sync.Mutex
	filename string
	samples  []loadcell.ForceSample
	closed   bool
}

func (u *csvFileRecorderUpdater) Write(samples ...loadcell.ForceSample) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.closed {
		return isometric.ErrWriteAfterClosed
	}
	u.samples = append(u.samples, samples...)
	return nil
}

func (u *csvFileRecorderUpdater) Finish(_ isometric.WorkoutOutcome) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.closed {
		return isometric.ErrWriteAfterClosed
	}

	if len(u.samples) == 0 {
		return isometric.ErrNoData
	}

	f, err := os.OpenFile(u.filename, os.O_RDWR|os.O_CREATE, 0444)
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)

	columnHeaders := []string{
		"time", "force",
	}

	// write column headers
	if err := w.Write(columnHeaders); err != nil {
		return err
	}

	// ensure data is sorted to clean up our output
	sort.Slice(u.samples, func(i, j int) bool {
		return u.samples[i].Time.Before(u.samples[j].Time)
	})

	// we mainly care about duration from start
	start := u.samples[0].Time

	// write rest of data
	for _, s := range u.samples {
		entry := []string{
			fmt.Sprintf("%d", s.Sub(start)), fmt.Sprintf("%d", s.Force/physic.Newton),
		}
		if err := w.Write(entry); err != nil {
			return err
		}
	}
	w.Flush()
	u.close()
	return f.Close()
}

func (u *csvFileRecorderUpdater) Close() {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.close()
}

func (u *csvFileRecorderUpdater) close() {
	u.closed = true
}

func CsvRecorder(dir string) (isometric.WorkoutRecorder, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return &csvFileRecorder{dir: dir}, nil
}

func (r *csvFileRecorder) Start(_ context.Context, descriptor string) (isometric.WorkoutUpdater, error) {
	const timeFormat = "20060102150405"
	filename := fmt.Sprintf("%s-%s.csv", time.Now().Format(timeFormat), descriptor)
	fpath := filepath.Join(r.dir, filename)
	return &csvFileRecorderUpdater{
		filename: fpath,
	}, nil
}
