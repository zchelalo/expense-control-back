package domain

type Amount struct{ value float64 }

func NewAmount(v float64) (Amount, error) {
	if v < 0 {
		return Amount{}, ErrInvalidAmount
	}
	return Amount{value: v}, nil
}

func (a Amount) Float64() float64 { return a.value }
