package testhang

import (
	"time"

	"github.com/spf13/cobra"
)

const (
	flagDuration = "duration"
)

func flags(cmd *cobra.Command) error {
	cmd.Flags().DurationP(flagDuration, "d", 12*time.Second, "time interval for max hang test")
	return nil
}
