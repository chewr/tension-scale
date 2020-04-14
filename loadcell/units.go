package loadcell

type ForceUnit int32

const (
	UNIT_RAW ForceUnit = 1
	UNIT_N   ForceUnit = 737  // g ~= 9.79945 m/s^2 in MPK
	UNIT_LB  ForceUnit = 3276 // 2.205 lbs per kg
	UNIT_KG  ForceUnit = 7222 // measured 7217 ~ 7227
)

type Newtons int
type Pounds int
type Kilograms int

func RawToNewtons(raw int32) Pounds {
	return Pounds(raw / int32(UNIT_N))
}

func ConvertToPounds(raw int32) Pounds {
	return Pounds(raw / int32(UNIT_LB))
}

func ConvertToKilograms(raw int32) Kilograms {
	return Kilograms(raw / int32(UNIT_KG))
}

func PoundsToRaw(lb Pounds) ForceUnit {
	return ForceUnit(UNIT_LB * ForceUnit(lb))
}

func KilogramsToRaw(kg Kilograms) ForceUnit {
	return ForceUnit(UNIT_KG * ForceUnit(kg))
}

func NewtonsToRaw(n Newtons) ForceUnit {
	return ForceUnit(UNIT_N * ForceUnit(n))
}
