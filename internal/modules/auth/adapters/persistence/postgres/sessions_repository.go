package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	authdb "github.com/zchelalo/expense-control-back/internal/db/sqlc/auth"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/domain"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/ports"
)

type SessionRepo struct {
  q *authdb.Queries
}

func NewSessionRepo(db authdb.DBTX) *SessionRepo {
  return &SessionRepo{q: authdb.New(db)}
}

func (r *SessionRepo) Create(ctx context.Context, s domain.Session) error {
	var remokedAt pgtype.Timestamptz
	if s.RevokedAt() != nil {
		remokedAt = pgtype.Timestamptz{Time: *s.RevokedAt(), Valid: true}
	} else {
		remokedAt = pgtype.Timestamptz{Valid: false}
	}
	params := authdb.CreateSessionParams{
		ID: pgtype.UUID{
			Bytes: s.ID().UUID(),
			Valid: true,
		},
		UserID: pgtype.UUID{
			Bytes: s.SubjectID().UUID(),
			Valid: true,
		},
		RefreshJti: pgtype.UUID{
			Bytes: s.RefreshID().UUID(),
			Valid: true,
		},
		ExpiresAt: pgtype.Timestamptz{
			Time: s.ExpiresAt(),
			Valid: true,
		},
		CreatedAt: pgtype.Timestamptz{
			Time: s.CreatedAt(),
			Valid: true,
		},
		RevokedAt: remokedAt,
	}

	err := r.q.CreateSession(ctx, params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "auth_sessions_refresh_jti_key" {
				return ports.ErrAlreadyExists{Name: "refresh jti"}
			}
		}
		return err
	}

	return nil
}

func (r *SessionRepo) ByID(ctx context.Context, id domain.SessionID) (domain.Session, error) {
	session, err := r.q.GetSessionByID(ctx, pgtype.UUID{
		Bytes: id.UUID(),
		Valid: true,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Session{}, ports.ErrNotFound{Name:"session"}
		}
		return domain.Session{}, err
	}

	parsedID, err := domain.NewSessionID(session.ID.Bytes)
	if err != nil {
		return domain.Session{}, err
	}
	parsedSubID, err := domain.NewSubjectID(session.UserID.Bytes)
	if err != nil {
		return domain.Session{}, err
	}
	parsedRefreshID, err := domain.NewRefreshTokenID(session.RefreshJti.Bytes)
	if err != nil {
		return domain.Session{}, err
	}
	var parsedRevokedAt *time.Time
	if session.RevokedAt.Valid {
		t := session.RevokedAt.Time
		parsedRevokedAt = &t
	}

	return domain.RehydrateSession(
		parsedID,
		parsedSubID,
		parsedRefreshID,
		session.CreatedAt.Time,
		session.ExpiresAt.Time,
		parsedRevokedAt,
	), nil
}

func (r *SessionRepo) RotateRefresh(ctx context.Context, sessionID domain.SessionID, newRefreshID domain.RefreshTokenID, newExp time.Time) error {
	err := r.q.RotateSessionRefreshID(ctx, authdb.RotateSessionRefreshIDParams{
		ID: pgtype.UUID{
			Bytes: sessionID.UUID(),
			Valid: true,
		},
		RefreshJti: pgtype.UUID{
			Bytes: newRefreshID.UUID(),
			Valid: true,
		},
		ExpiresAt: pgtype.Timestamptz{
			Time: newExp,
			Valid: true,
		},
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "auth_sessions_refresh_jti_key" {
				return ports.ErrAlreadyExists{Name: "refresh jti"}
			}
		}
		return err
	}

	return nil
}

func (r *SessionRepo) Revoke(ctx context.Context, sessionID domain.SessionID, now time.Time) error {
	err := r.q.RevokeSession(ctx, authdb.RevokeSessionParams{
		ID: pgtype.UUID{
			Bytes: sessionID.UUID(),
			Valid: true,
		},
		RevokedAt: pgtype.Timestamptz{
			Time: now,
			Valid: true,
		},
	})
	if err != nil {
		return err
	}

	return nil
}