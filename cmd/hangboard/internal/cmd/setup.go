package cmd

import (
	"github.com/chewr/tension-scale/cmd/hangboard/internal/cmd/daemon"
	"github.com/chewr/tension-scale/cmd/hangboard/internal/cmd/workout"
	"github.com/spf13/cobra"
)

func setup(rootCmd *cobra.Command) error {
	// add flags...
	workout.AddCommands(rootCmd)
	daemon.AddCommands(rootCmd)
	return nil
}
