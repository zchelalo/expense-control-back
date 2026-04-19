package domain

type Description struct{ value string }

func NewDescription(v string) (Description, error) {
	if v == "" {
		return Description{}, ErrInvalidDescription
	}
	return Description{value: v}, nil
}

func (d Description) String() string { return d.value }
