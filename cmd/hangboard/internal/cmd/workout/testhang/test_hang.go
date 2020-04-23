package testhang

import (
	"time"

	"github.com/chewr/tension-scale/cmd/hangboard/internal/cmd/workout/shared"
	"github.com/chewr/tension-scale/isometric"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run a hangboard test",
	Long:  `Test the max pull for a specified duration`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: doMaxTest,
}

func AddCommands(rootCmd *cobra.Command) {
	flags(testCmd)
	rootCmd.AddCommand(testCmd)
}

func doMaxTest(cmd *cobra.Command, args []string) error {
	display, err := shared.SetupDisplay()
	if err != nil {
		return err
	}
	loadCell, err := shared.SetupLoadCell()
	if err != nil {
		return err
	}
	if err := loadCell.Tare(cmd.Context(), 20); err != nil {
		return err
	}
	recorder, err := shared.SetupOutput()
	if err != nil {
		return err
	}

	duration, err := cmd.Flags().GetDuration(flagDuration)
	if err != nil {
		return err
	}

	maxTestWorkout := setupMaxTestWorkout(duration)
	return maxTestWorkout.Run(cmd.Context(), display, loadCell, recorder)
}

func setupMaxTestWorkout(d time.Duration) isometric.Workout {
	return isometric.Composite(
		isometric.SetupInterval(),
		isometric.RestInterval(5*time.Second),
		isometric.MaxTest(d),
		isometric.RestInterval(time.Second*5),
	)
}
