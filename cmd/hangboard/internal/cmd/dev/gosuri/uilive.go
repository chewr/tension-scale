package gosuri

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
	"github.com/gosuri/uilive"
	"github.com/spf13/cobra"
)

var spinnerCmd = &cobra.Command{
	Use:   "uilive",
	Short: "test the spinner",
	RunE:  doUilive,
}

func AddCommands(rootCmd *cobra.Command) {
	addFlags(spinnerCmd)
	rootCmd.AddCommand(spinnerCmd)
}

const (
	flagColor = "color"
)

func addFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(flagColor, "c", "", "color")
}

func doUilive(cmd *cobra.Command, args []string) error {
	strings := []string{
		"foo",
		"bar",
		"foobar",
	}
	colors := [][]color.Attribute{
		{color.BlinkRapid, color.FgCyan},
		{color.Italic},
		{color.FgHiRed},
		{color.BgBlue, color.FgMagenta},
	}
	uw := uilive.New()
	uw.Out = cmd.OutOrStdout()
	for {
		select {
		case <-cmd.Context().Done():
			return nil
		default:
		}
		s := strings[rand.Intn(len(strings))]
		c := colors[rand.Intn(len(colors))]

		if _, err := fmt.Fprintln(uw, color.New(c...).SprintfFunc()(s)); err != nil {
			return err
		}
		if err := uw.Flush(); err != nil {
			return err
		}
		time.Sleep(time.Second / 2)
	}
}
