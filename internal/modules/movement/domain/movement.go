package domain

import "time"

type Movement struct {
	id             MovementID
	amount         Amount
	description    Description
	movementTypeID MovementTypeID
	categoryID     CategoryID
	accountID      AccountID
	userID         UserID
	createdAt      time.Time
	updatedAt      time.Time
	deletedAt      *time.Time
}

func NewMovement(
	id MovementID,
	amount Amount,
	description Description,
	movementTypeID MovementTypeID,
	categoryID CategoryID,
	accountID AccountID,
	userID UserID,
	now time.Time,
) Movement {
	return Movement{
		id:             id,
		amount:         amount,
		description:    description,
		movementTypeID: movementTypeID,
		categoryID:     categoryID,
		accountID:      accountID,
		userID:         userID,
		createdAt:      now,
		updatedAt:      now,
		deletedAt:      nil,
	}
}

func RehydrateMovement(
	id MovementID,
	amount Amount,
	description Description,
	movementTypeID MovementTypeID,
	categoryID CategoryID,
	accountID AccountID,
	userID UserID,
	createdAt, updatedAt time.Time,
	deletedAt *time.Time,
) Movement {
	return Movement{
		id:             id,
		amount:         amount,
		description:    description,
		movementTypeID: movementTypeID,
		categoryID:     categoryID,
		accountID:      accountID,
		userID:         userID,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
		deletedAt:      deletedAt,
	}
}

func (m Movement) ID() MovementID                 { return m.id }
func (m Movement) Amount() Amount                 { return m.amount }
func (m Movement) Description() Description       { return m.description }
func (m Movement) MovementTypeID() MovementTypeID { return m.movementTypeID }
func (m Movement) CategoryID() CategoryID         { return m.categoryID }
func (m Movement) AccountID() AccountID           { return m.accountID }
func (m Movement) UserID() UserID                 { return m.userID }
func (m Movement) CreatedAt() time.Time           { return m.createdAt }
func (m Movement) UpdatedAt() time.Time           { return m.updatedAt }
func (m Movement) DeletedAt() *time.Time          { return m.deletedAt }
func (m Movement) IsDeleted() bool                { return m.deletedAt != nil }

func (m *Movement) SoftDelete(now time.Time) {
	m.deletedAt = &now
	m.updatedAt = now
}
