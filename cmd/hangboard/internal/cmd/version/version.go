package version

import (
	"fmt"

	"github.com/chewr/tension-scale/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), version.GetVersion()); err != nil {
			return err
		}
		return nil
	},
}

func AddCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(versionCmd)
}
