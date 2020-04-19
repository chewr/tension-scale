package maxhang

import (
	"github.com/chewr/tension-scale/cmd/hangboard/internal/wip"
	"github.com/spf13/cobra"
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
	RunE: func(cmd *cobra.Command, args []string) error { return wip.ErrCmdNotImplemented },
}

func AddCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(maxHangCmd)
}
