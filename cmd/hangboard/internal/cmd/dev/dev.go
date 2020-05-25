package dev

import (
	"github.com/chewr/tension-scale/cmd/hangboard/internal/cmd/dev/refresh"
	"github.com/chewr/tension-scale/cmd/hangboard/internal/cmd/dev/spintest"
	"github.com/spf13/cobra"
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Dev tools",
}

func AddCommands(rootCmd *cobra.Command) {
	setup(devCmd)
	rootCmd.AddCommand(devCmd)
}

func setup(devCmd *cobra.Command) {
	spintest.AddCommands(devCmd)
	refresh.AddCommands(devCmd)
}
