package maxhang

import (
	"github.com/spf13/cobra"
)

const (
	flagThreshold = "threshold"
	flagWeek      = "week"
)

func flags(cmd *cobra.Command) error {
	cmd.Flags().StringP(flagThreshold, "t", "0N", "force threshold for workout")
	if err := cmd.MarkFlagRequired(flagThreshold); err != nil {
		return err
	}
	cmd.Flags().IntP(flagWeek, "w", 1, "week for max hang workout")
	if err := cmd.MarkFlagRequired(flagThreshold); err != nil {
		return err
	}
	return nil
}
