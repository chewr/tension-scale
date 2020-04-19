package workout

import (
	"github.com/chewr/tension-scale/cmd/hangboard/internal/cmd/workout/maxhang"
	"github.com/chewr/tension-scale/cmd/hangboard/internal/cmd/workout/testhang"
	"github.com/spf13/cobra"
)

var workoutCmd = &cobra.Command{
	Use:   "hangboard",
	Short: "Run a workout with the hangboard",
}

func setup(workoutCmd *cobra.Command) {
	// add flags...
	maxhang.AddCommands(workoutCmd)
	testhang.AddCommands(workoutCmd)
}

func AddCommands(rootCmd *cobra.Command) {
	setup(workoutCmd)
	rootCmd.AddCommand(workoutCmd)
}
