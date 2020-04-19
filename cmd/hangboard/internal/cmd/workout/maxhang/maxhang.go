package maxhang

import (
	"github.com/chewr/tension-scale/workout/maxhang"
	"github.com/spf13/cobra"
	"periph.io/x/periph/conn/physic"
)

var maxHangCmd = &cobra.Command{
	Use:   "max-hang",
	Short: "Run a max hang workout",
	Long: `Max hang workouts are designed to increase maximum force
output. A max hang protocol involves exerting maximum or 
near-maximum effort for a short amount of time. A typical
program consists of a four-week cycle with each workout
performed 2-3 times per week:

    Week 1:
        3", 6", 9" work intervals on 30" centers
        3 cycles with 90" rest
    Week 2:
        3", 6", 9" work intervals on 30" centers
        4 cycles with 90" rest
    Week 3:
        3", 6", 9" work intervals on 30" centers
        5 cycles with 90" rest
    Week 4:
        3", 6", 9", 12" work intervals on 30" centers
        3 cycles with 90" rest

For more information, see 
https://www.climbstrong.com/education-center/making-sense-hangboard-programs/
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: doWorkout,
}

func AddCommands(rootCmd *cobra.Command) {
	flags(maxHangCmd)
	rootCmd.AddCommand(maxHangCmd)
}

func doWorkout(cmd *cobra.Command, args []string) error {
	threshold, err := cmd.Flags().GetInt64(flagThreshold)
	if err != nil {
		return err
	}
	week, err := cmd.Flags().GetInt(flagWeek)
	if err != nil {
		return err
	}
	workout, err := maxhang.MaxHangWorkout(maxhang.Week(week), physic.Force(threshold)*physic.PoundForce)
	if err != nil {
		return err
	}
	display, err := SetupDisplay()
	if err != nil {
		return err
	}
	loadCell, err := SetupLoadCell()
	if err != nil {
		return err
	}
	return workout.Run(cmd.Context(), display, loadCell)
}
