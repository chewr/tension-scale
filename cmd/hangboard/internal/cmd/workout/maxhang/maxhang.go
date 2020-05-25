package maxhang

import (
	"github.com/chewr/tension-scale/cmd/hangboard/internal/cmd/workout/recording"
	"github.com/chewr/tension-scale/cmd/hangboard/internal/cmd/workout/shared"
	"github.com/chewr/tension-scale/display"
	"github.com/chewr/tension-scale/display/cli"
	"github.com/chewr/tension-scale/errutil"
	"github.com/chewr/tension-scale/isometric/data"
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
	errutil.PanicOnErr(flags(maxHangCmd))
	rootCmd.AddCommand(maxHangCmd)
}

func doWorkout(cmd *cobra.Command, args []string) error {
	threshold, err := cmd.Flags().GetString(flagThreshold)
	if err != nil {
		return err
	}
	f := new(physic.Force)
	if err := f.Set(threshold); err != nil {
		return err
	}
	week, err := cmd.Flags().GetInt(flagWeek)
	if err != nil {
		return err
	}
	maxHangWorkout, err := maxhang.Workout(maxhang.Week(week), *f)
	if err != nil {
		return err
	}
	ledDisplay, err := shared.SetupDisplay()
	if err != nil {
		return err
	}
	ledDisplay.Start(cmd.Context())
	loadCell, err := shared.SetupLoadCell()
	if err != nil {
		return err
	}
	if err := loadCell.Tare(cmd.Context(), 20); err != nil {
		return err
	}
	fileRecorder, err := shared.SetupOutput()
	if err != nil {
		return err
	}

	// TODO(rchew) reconcile cliRecorder with the cliDisplay
	cliRecorder := recording.CliRecorder(cmd)
	recorder := data.MultiRecorder(fileRecorder, cliRecorder)

	cliModel, err := cli.NewCliDisplay(cmd.OutOrStdout())
	if err != nil {
		return err
	}
	cliModel.Start(cmd.Context())

	model := display.ModelMux(ledDisplay, cliModel)

	return maxHangWorkout.Run(cmd.Context(), model, loadCell, recorder)
}
