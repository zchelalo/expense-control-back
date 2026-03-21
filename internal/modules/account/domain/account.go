package domain

import "time"

type Account struct {
	id        AccountID
	name      Name
	balance   Balance
	userId    UserID
	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

func NewAccount(id AccountID, name Name, balance Balance, userId UserID, now time.Time) Account {
	return Account{
		id:        id,
		name:      name,
		balance:   balance,
		userId:    userId,
		createdAt: now,
		updatedAt: now,
		deletedAt: nil,
	}
}

func RehydrateAccount(
	id AccountID,
	name Name,
	balance Balance,
	userId UserID,
	createdAt, updatedAt time.Time,
	deletedAt *time.Time,
) Account {
	return Account{
		id:        id,
		name:     name,
		balance:  balance,
		userId:  userId,
		createdAt: createdAt,
		updatedAt: updatedAt,
		deletedAt: deletedAt,
	}
}

func (a Account) ID() AccountID            { return a.id }
func (a Account) Name() Name               { return a.name }
func (a Account) Balance() Balance         { return a.balance }
func (a Account) UserID() UserID           { return a.userId }
func (a Account) CreatedAt() time.Time     { return a.createdAt }
func (a Account) UpdatedAt() time.Time     { return a.updatedAt }
func (a Account) DeletedAt() *time.Time    { return a.deletedAt }
func (a Account) IsDeleted() bool          { return a.deletedAt != nil }

func (a *Account) SoftDelete(now time.Time) {
	a.deletedAt = &now
	a.updatedAt = now
}