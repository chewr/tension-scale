package dev

import (
	"github.com/chewr/tension-scale/cmd/hangboard/internal/wip"
	"github.com/spf13/cobra"
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Dev tools",
	RunE:  func(cmd *cobra.Command, args []string) error { return wip.ErrTODO },
}

func AddCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(devCmd)
}
