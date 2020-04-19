package loadcell

import (
	"periph.io/x/periph/conn/physic"
)

var (
	// TrueSun400Slow represents measured values for the
	// TrueSun 400KG-capacity loadcell measured at 10SPS
	// We determined that 1KG ~= 7222 units
	TrueSun400Slow = Calibrate(7222, physic.EarthGravity)
)

// Calibration implements methods for converting from measured
// values to real force units
type Calibration interface {
	ToForce(int64) physic.Force
}

// insensitiveCalibration is intended to calibrate a device for which
// the minimum reading is greater than the minimum force that can be
// expressed
type insensitiveCalibration int64

func (c insensitiveCalibration) ToForce(raw int64) physic.Force {
	return physic.Force(raw * int64(c))
}

// sensitiveCalibration is intended to calibrate a device for which
// the minimum reading is less than the minimum force that can be
// expressed
type sensitiveCalibration int64

func (c sensitiveCalibration) ToForce(raw int64) physic.Force {
	return physic.Force(raw / int64(c))
}

// amplifiedCalibration is intended to calibrate a device for which
// the minimum reading is within several orders of magnitude of the
// minimum force that can be expressed. This is to avoid magnifying
// integer resolution errors, but reduces the maximum force which
// can be read.
type amplifiedCalibration struct {
	amp, conversion int64
}

func (c amplifiedCalibration) ToForce(raw int64) physic.Force {
	amplified := raw * c.amp
	return physic.Force(amplified / c.conversion)
}

func Calibrate(reading int64, actual physic.Force) Calibration {
	const minGain = 1000
	if actual/physic.Force(reading) > minGain {
		return insensitiveCalibration(actual / physic.Force(reading))
	}
	if physic.Force(reading)/actual > minGain {
		return sensitiveCalibration(physic.Force(reading) / actual)
	}
	amplified := reading * minGain
	return amplifiedCalibration{
		amp:        minGain,
		conversion: amplified / int64(actual),
	}
}
