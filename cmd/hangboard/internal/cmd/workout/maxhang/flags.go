package maxhang

import (
	"github.com/spf13/cobra"
)

const (
	flagThreshold = "threshold"
	flagWeek      = "week"
)

func flags(cmd *cobra.Command) error {
	cmd.Flags().Int64P(flagThreshold, "t", 0, "weigh (pounds) for workout")
	if err := cmd.MarkFlagRequired(flagThreshold); err != nil {
		return err
	}
	cmd.Flags().IntP(flagWeek, "w", 1, "week for max hang workout")
	if err := cmd.MarkFlagRequired(flagThreshold); err != nil {
		return err
	}
	return nil
}
