package parse

import "github.com/google/uuid"

func UUID(raw string) (uuid.UUID, error) {
	return uuid.Parse(raw)
}

func OptionalUUID(raw string) (*uuid.UUID, error) {
	if raw == "" {
		return nil, nil
	}

	id, err := uuid.Parse(raw)
	if err != nil {
		return nil, err
	}

	return &id, nil
}
