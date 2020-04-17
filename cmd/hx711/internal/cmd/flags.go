package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagContinuous              = "continuous"
	flagDebug                   = "debug"
	flagGain                    = "gain"
	flagInstantaneousRead       = "instantaneous"
	flagReset                   = "reset"
	flagSamples                 = "samples"
	flagUsePeriphImplementation = "use-periph-implementation"
)

func setupCmd() error {
	readCmd.Flags().IntP(flagSamples, "s", -1, "read out n samples; -1 for continuous")
	if err := viper.BindPFlag(flagSamples, readCmd.Flag(flagSamples)); err != nil {
		return err
	}
	readCmd.Flags().BoolP(flagUsePeriphImplementation, "x", false, "use implementation of hx711 from periph.io")
	readCmd.Flag(flagUsePeriphImplementation).NoOptDefVal = "true"
	if err := viper.BindPFlag(flagUsePeriphImplementation, readCmd.Flag(flagUsePeriphImplementation)); err != nil {
		return err
	}
	readCmd.Flags().BoolP(flagReset, "r", false, "reset hardware on startup (only applies when --use-periph-implementation is false)")
	readCmd.Flag(flagReset).NoOptDefVal = "true"
	if err := viper.BindPFlag(flagReset, readCmd.Flag(flagReset)); err != nil {
		return err
	}
	readCmd.Flags().BoolP(flagContinuous, "c", false, "read values using continuous implementation")
	readCmd.Flag(flagContinuous).NoOptDefVal = "true"
	if err := viper.BindPFlag(flagContinuous, readCmd.Flag(flagContinuous)); err != nil {
		return err
	}
	readCmd.Flags().BoolP(flagInstantaneousRead, "i", false, "force reads to be instantaneous. Only valid when not reading continuously")
	readCmd.Flag(flagInstantaneousRead).NoOptDefVal = "true"
	if err := viper.BindPFlag(flagInstantaneousRead, readCmd.Flag(flagInstantaneousRead)); err != nil {
		return err
	}
	readCmd.Flags().BoolP(flagDebug, "d", false, "enable debug logging")
	readCmd.Flag(flagDebug).NoOptDefVal = "true"
	if err := viper.BindPFlag(flagDebug, readCmd.Flag(flagDebug)); err != nil {
		return err
	}
	readCmd.Flags().IntP(flagGain, "g", 128, "set gain (valid values: 32, 64, 128")
	if err := viper.BindPFlag(flagGain, readCmd.Flag(flagGain)); err != nil {
		return err
	}
	return nil
}

func validateFlags(cmd *cobra.Command) {
	if viper.GetBool(flagUsePeriphImplementation) && viper.GetBool(flagReset) {
		cmd.PrintErrln("reset only applies when use-periph-implementation is unset, will ignore reset flag")
		viper.Set(flagReset, false)
	}
	if viper.GetBool(flagInstantaneousRead) && viper.GetBool(flagContinuous) {
		cmd.PrintErrln("instantaneous flag and continuous read modes are mutually exclusive, will ignore instantaneous flag")
		viper.Set(flagInstantaneousRead, false)
	}
	switch gain := viper.GetInt(flagGain); gain {
	case 128:
	case 32:
	case 64:
	default:
		cmd.PrintErrln(fmt.Sprintf("invalid gain of %d; defaulting to 128", gain))
		viper.Set(flagGain, 128)
	}
}
