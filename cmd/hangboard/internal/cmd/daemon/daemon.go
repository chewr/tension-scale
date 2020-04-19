package daemon

import (
	"github.com/chewr/tension-scale/cmd/hangboard/internal/wip"
	"github.com/spf13/cobra"
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run a hangboard daemon",
	RunE:  func(cmd *cobra.Command, args []string) error { return wip.ErrTODO },
}
