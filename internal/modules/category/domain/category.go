package domain

import "time"

type Category struct {
	id        CategoryID
	name      Name
	userID    UserID
	createdAt time.Time
	updatedAt time.Time
}

func NewCategory(
	id CategoryID,
	name Name,
	userID UserID,
	now time.Time,
) Category {
	return Category{
		id:        id,
		name:      name,
		userID:    userID,
		createdAt: now,
		updatedAt: now,
	}
}

func RehydrateCategory(
	id CategoryID,
	rawName string,
	userID UserID,
	createdAt time.Time,
	updatedAt time.Time,
) (Category, error) {
	name, err := NewName(rawName)
	if err != nil {
		return Category{}, err
	}

	return Category{
		id:        id,
		name:      name,
		userID:    userID,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

func (c Category) ID() CategoryID       { return c.id }
func (c Category) Name() Name           { return c.name }
func (c Category) UserID() UserID       { return c.userID }
func (c Category) CreatedAt() time.Time { return c.createdAt }
func (c Category) UpdatedAt() time.Time { return c.updatedAt }
