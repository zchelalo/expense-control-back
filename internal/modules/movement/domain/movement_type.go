package domain

const (
	MovementTypeKeyIncome  = "income"
	MovementTypeKeyExpense = "expense"
)

type MovementType struct {
	id   MovementTypeID
	key  string
	name string
}

func RehydrateMovementType(id MovementTypeID, key, name string) (MovementType, error) {
	if key == "" || name == "" {
		return MovementType{}, ErrInvalidMovementType
	}

	return MovementType{
		id:   id,
		key:  key,
		name: name,
	}, nil
}

func (mt MovementType) ID() MovementTypeID { return mt.id }
func (mt MovementType) Key() string        { return mt.key }
func (mt MovementType) Name() string       { return mt.name }
