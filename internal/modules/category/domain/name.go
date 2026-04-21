package domain

type Name struct{ value string }

func NewName(v string) (Name, error) {
	if v == "" {
		return Name{}, ErrInvalidName
	}

	return Name{value: v}, nil
}

func (n Name) String() string { return n.value }
