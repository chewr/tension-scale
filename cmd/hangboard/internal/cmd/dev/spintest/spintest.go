package spintest

import (
	"context"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var spinnerCmd = &cobra.Command{
	Use:   "spinner",
	Short: "test the spinner",
	RunE:  doSpinner,
}

func AddCommands(rootCmd *cobra.Command) {
	addFlags(spinnerCmd)
	rootCmd.AddCommand(spinnerCmd)
}

const (
	flagSuffix         = "suffix"
	flagDefaultCharset = "default-charset"
)

func addFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(flagSuffix, "s", "", "spinner suffix")
	cmd.Flags().IntP(flagDefaultCharset, "d", 0, "default charset")
}

func doSpinner(cmd *cobra.Command, args []string) error {
	var opts []spinner.Option
	opts = append(opts, spinner.WithWriter(cmd.OutOrStdout()))

	suffix, err := cmd.Flags().GetString(flagSuffix)
	if err != nil {
		return err
	}
	if suffix != "" {
		opts = append(opts, spinner.WithSuffix(suffix))
	}

	charSetIndex, err := cmd.Flags().GetInt(flagDefaultCharset)
	if err != nil {
		return err
	}

	s := spinner.New(spinner.CharSets[charSetIndex], 300*time.Millisecond, opts...)
	s.Start()
	defer s.Stop()

	ctx := cmd.Context()
	<-ctx.Done()
	if err := ctx.Err(); err != nil && err != context.Canceled {
		return err
	}
	return nil
}
