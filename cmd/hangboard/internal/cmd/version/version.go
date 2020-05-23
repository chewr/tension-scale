package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "none"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), version); err != nil {
			return err
		}
		return nil
	},
}

func AddCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(versionCmd)
}
