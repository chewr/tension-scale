package testhang

import (
	"github.com/chewr/tension-scale/cmd/hangboard/internal/wip"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run a hangboard test",
	Long:  `todo`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error { return wip.ErrTODO },
}

func AddCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(testCmd)
}
