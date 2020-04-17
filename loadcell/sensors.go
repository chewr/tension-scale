package loadcell

var (
	// TrueSun400 represents measured values for the
	// TrueSun 400KG-capacity loadcell
	TrueSun400 = calibrateByKilograms(7222) // measurement of 7217 ~ 7227
)

// Calibration implements methods for converting from measured
// values to real force units
type Calibration interface {
	Newtons(int32) Newtons
	Pounds(int32) Pounds
	Kilograms(int32) Kilograms
}

var _ Calibration = new(kilogramCalibrated)

// kilogramCalibrated is an implementation of Calibration
// based on the raw unit reading for a known mass
type kilogramCalibrated int32

func (c kilogramCalibrated) Kilograms(raw int32) Kilograms { return Kilograms(raw * int32(c)) }
func (c kilogramCalibrated) Pounds(raw int32) Pounds       { return PoundsFromKilograms(c.Kilograms(raw)) }
func (c kilogramCalibrated) Newtons(raw int32) Newtons     { return NewtonsFromKilograms(c.Kilograms(raw)) }

func calibrateByKilograms(factor int32) Calibration {
	return kilogramCalibrated(factor)
}
