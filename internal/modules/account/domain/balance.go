package domain

type Balance struct{ value float64 }

func NewBalance(v float64) (Balance, error) {
	if v < 0 {
		return Balance{}, ErrInvalidBalance
	}
	return Balance{value: v}, nil
}

func (b Balance) Float64() float64 { return b.value }