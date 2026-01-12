package domain

import "time"

type Account struct {
	id        SubjectID
	email     Email
	passHash  PasswordHash
	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

func NewAccount(id SubjectID, email Email, passHash PasswordHash, now time.Time) Account {
	return Account{
		id:        id,
		email:     email,
		passHash:  passHash,
		createdAt: now,
		updatedAt: now,
		deletedAt: nil,
	}
}

func RehydrateAccount(
	id SubjectID,
	email Email,
	passHash PasswordHash,
	createdAt, updatedAt time.Time,
	deletedAt *time.Time,
) Account {
	return Account{
		id:        id,
		email:     email,
		passHash:  passHash,
		createdAt: createdAt,
		updatedAt: updatedAt,
		deletedAt: deletedAt,
	}
}

func (a Account) ID() SubjectID            { return a.id }
func (a Account) Email() Email             { return a.email }
func (a Account) PasswordHash() PasswordHash { return a.passHash }
func (a Account) CreatedAt() time.Time     { return a.createdAt }
func (a Account) UpdatedAt() time.Time     { return a.updatedAt }
func (a Account) DeletedAt() *time.Time    { return a.deletedAt }
func (a Account) IsDeleted() bool          { return a.deletedAt != nil }

func (a Account) CanAuthenticate() bool {
	return !a.IsDeleted()
}

func (a *Account) SoftDelete(now time.Time) {
	a.deletedAt = &now
	a.updatedAt = now
}