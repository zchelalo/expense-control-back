package domain

import "time"

type User struct {
	id        UserID
	email     Email
	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

func NewUser(id UserID, email Email, now time.Time) User {
	return User{
		id:        id,
		email:     email,
	}
}

func RehydrateUser(
	id UserID,
	email Email,
	createdAt, updatedAt time.Time,
	deletedAt *time.Time,
) User {
	return User{
		id:        id,
		email:     email,
		createdAt: createdAt,
		updatedAt: updatedAt,
		deletedAt: deletedAt,
	}
}

func (u User) ID() UserID             { return u.id }
func (u User) Email() Email           { return u.email }
func (u User) CreatedAt() time.Time   { return u.createdAt }
func (u User) UpdatedAt() time.Time   { return u.updatedAt }
func (u User) DeletedAt() *time.Time  { return u.deletedAt }
func (u User) IsDeleted() bool        { return u.deletedAt != nil }

func (u *User) MarkUpdated(now time.Time) { u.updatedAt = now }

func (u *User) SoftDelete(now time.Time) {
	u.deletedAt = &now
	u.updatedAt = now
}