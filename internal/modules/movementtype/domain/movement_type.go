package domain

type MovementType struct {
	id          MovementTypeID
	key         string
	name        string
	description string
}

func RehydrateMovementType(id MovementTypeID, key, name, description string) (MovementType, error) {
	if key == "" || name == "" {
		return MovementType{}, ErrInvalidMovementType
	}

	return MovementType{
		id:          id,
		key:         key,
		name:        name,
		description: description,
	}, nil
}

func (mt MovementType) ID() MovementTypeID { return mt.id }
func (mt MovementType) Key() string        { return mt.key }
func (mt MovementType) Name() string       { return mt.name }
func (mt MovementType) Description() string {
	return mt.description
}
