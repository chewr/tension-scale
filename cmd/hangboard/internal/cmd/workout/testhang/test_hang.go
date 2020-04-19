package testhang

import (
	"github.com/spf13/cobra"
	"github.com/chewr/tension-scale/cmd/hangboard/internal/wip"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run a hangboard test",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error { return wip.ErrCmdNotImplemented },
}

func AddCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(testCmd)
}
