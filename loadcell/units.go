package loadcell

type Newtons float32
type Pounds float32
type Kilograms float32

const (
	newtonsPerKilogram float32 = 9.799   // m/s^2 (gravitational acceleration)
	poundsPerKilogram  float32 = 2.20462 // lb/kg (defined unitless constant)
	newtonsPerPound    float32 = 4.4448  // n/lb  (gravitational acceleration, in weird units)
)

func NewtonsFromKilograms(kg Kilograms) Newtons { return Newtons(float32(kg) * newtonsPerKilogram) }
func NewtonsFromPounds(lb Pounds) Newtons       { return Newtons(float32(lb) * newtonsPerPound) }

func KilogramsFromPounds(lb Pounds) Kilograms  { return Kilograms(float32(lb) / poundsPerKilogram) }
func KilogramsFromNewtons(n Newtons) Kilograms { return Kilograms(float32(n) / newtonsPerKilogram) }

func PoundsFromKilograms(kg Kilograms) Pounds { return Pounds(float32(kg) * poundsPerKilogram) }
func PoundsFromNewtons(n Newtons) Pounds      { return Pounds(float32(n) / newtonsPerPound) }
