package testhang

import (
	"github.com/spf13/cobra"
	"time"
)

const (
	flagDuration = "duration"
)

func flags(cmd *cobra.Command) error {
	cmd.Flags().DurationP(flagDuration, "d", 12 * time.Second, "time interval for max hang test")
	return nil
}
