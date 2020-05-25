package refresh

import (
	"math/rand"
	"time"

	"github.com/chewr/tension-scale/display/cli"
	"github.com/chewr/tension-scale/display/cli/refresh"
	"github.com/chewr/tension-scale/display/input"
	"github.com/chewr/tension-scale/display/state"
	"github.com/spf13/cobra"
	"periph.io/x/periph/conn/physic"
)

var refreshingCliCmd = &cobra.Command{
	Use:   "refreshing-cli",
	Short: "test the refreshing cli",
	RunE:  doRefreshingCli,
}

func AddCommands(rootCmd *cobra.Command) {
	addFlags(refreshingCliCmd)
	rootCmd.AddCommand(refreshingCliCmd)
}

const (
	flagColor = "color"
)

func addFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(flagColor, "c", "", "color")
}

func doRefreshingCli(cmd *cobra.Command, args []string) error {
	// st := state.Rest(time.Now().Add(20 * time.Second))
	forceInput := &input.DynamicForceInput{}
	st := state.Work(
		input.ForceRequired(700*physic.Newton),
		forceInput,
		time.Now().Add(time.Second*12),
	)
	st.GetMutableState().Start()

	printer := refresh.NewPrinter(cmd.OutOrStdout())
	var f physic.Force = 0
	rand.Seed(time.Now().Unix())
	for {
		select {
		case <-cmd.Context().Done():
			return nil
		default:
		}

		var n physic.Force = 0
		norm60 := rand.Int63n(50) + rand.Int63n(50) + rand.Int63n(50) - 75
		if f > 800*physic.Newton {
			n = physic.Force(norm60) - 40
		} else if f > 725*physic.Newton {
			n = physic.Force(norm60) - 10
		} else if f > 700*physic.Newton {
			n = physic.Force(norm60) + 10
		} else if f > 650*physic.Newton {
			n = physic.Force(norm60) + 20
		} else {
			n = 30
		}
		f += n * physic.Newton

		forceInput.UpdateForceInput(f)

		if err := printer.Print(cli.ToCliOutput(st)); err != nil {
			return err
		}
		time.Sleep(time.Millisecond * 100)
	}
}
