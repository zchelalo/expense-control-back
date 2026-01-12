package domain

import "time"

type Session struct {
	id          SessionID
	subjectID   SubjectID
	refreshID   RefreshTokenID
	createdAt   time.Time
	expiresAt   time.Time
	revokedAt   *time.Time
}

func NewSession(id SessionID, subjectID SubjectID, refreshID RefreshTokenID, now time.Time, expiresAt time.Time) (Session, error) {
	if !expiresAt.After(now) {
		return Session{}, ErrSessionExpired
	}

	return Session{
		id:        id,
		subjectID: subjectID,
		refreshID: refreshID,
		createdAt: now,
		expiresAt: expiresAt,
		revokedAt: nil,
	}, nil
}

func RehydrateSession(
	id SessionID,
	subjectID SubjectID,
	refreshID RefreshTokenID,
	createdAt, expiresAt time.Time,
	revokedAt *time.Time,
) Session {
	return Session{
		id:        id,
		subjectID: subjectID,
		refreshID: refreshID,
		createdAt: createdAt,
		expiresAt: expiresAt,
		revokedAt: revokedAt,
	}
}

func (s Session) ID() SessionID          { return s.id }
func (s Session) SubjectID() SubjectID   { return s.subjectID }
func (s Session) RefreshID() RefreshTokenID { return s.refreshID }
func (s Session) CreatedAt() time.Time   { return s.createdAt }
func (s Session) ExpiresAt() time.Time   { return s.expiresAt }
func (s Session) RevokedAt() *time.Time  { return s.revokedAt }

func (s Session) IsRevoked() bool { return s.revokedAt != nil }

func (s Session) IsExpired(now time.Time) bool {
	return !s.expiresAt.After(now)
}

func (s Session) ValidateUsable(now time.Time) error {
	if s.IsRevoked() {
		return ErrSessionRevoked
	}
	if s.IsExpired(now) {
		return ErrSessionExpired
	}
	return nil
}

func (s *Session) Revoke(now time.Time) error {
	if s.IsRevoked() {
		return ErrSessionRevoked
	}
	s.revokedAt = &now
	return nil
}

func (s *Session) RotateRefresh(newRefreshID RefreshTokenID, now time.Time, newExpiresAt time.Time) error {
	if err := s.ValidateUsable(now); err != nil {
		return err
	}
	if !newExpiresAt.After(now) {
		return ErrSessionExpired
	}
	s.refreshID = newRefreshID
	s.expiresAt = newExpiresAt
	return nil
}