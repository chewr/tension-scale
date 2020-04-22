package loadcell

import (
	"periph.io/x/periph/host"
)

func init() {
	if _, err := host.Init(); err != nil {
		panic(err)
	}
}
