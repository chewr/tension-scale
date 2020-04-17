package loadcell

type Newtons float64
type Pounds float64
type Kilograms float64

const (
	newtonsPerKilogram float64 = 9.799   // m/s^2 (gravitational acceleration)
	poundsPerKilogram  float64 = 2.20462 // lb/kg (defined unitless constant)
	newtonsPerPound    float64 = 4.4448  // n/lb  (gravitational acceleration, in weird units)
)

func NewtonsFromKilograms(kg Kilograms) Newtons { return Newtons(float64(kg) * newtonsPerKilogram) }
func NewtonsFromPounds(lb Pounds) Newtons       { return Newtons(float64(lb) * newtonsPerPound) }

func KilogramsFromPounds(lb Pounds) Kilograms  { return Kilograms(float64(lb) / poundsPerKilogram) }
func KilogramsFromNewtons(n Newtons) Kilograms { return Kilograms(float64(n) / newtonsPerKilogram) }

func PoundsFromKilograms(kg Kilograms) Pounds { return Pounds(float64(kg) * poundsPerKilogram) }
func PoundsFromNewtons(n Newtons) Pounds      { return Pounds(float64(n) / newtonsPerPound) }
