package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	authdb "github.com/zchelalo/expense-control-back/internal/db/sqlc/auth"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/domain"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/ports"
	pgutil "github.com/zchelalo/expense-control-back/internal/shared/postgresutil"
)

type SessionRepo struct {
	q *authdb.Queries
}

func NewSessionRepo(db authdb.DBTX) *SessionRepo {
	return &SessionRepo{q: authdb.New(db)}
}

func (r *SessionRepo) Create(ctx context.Context, s domain.Session) error {
	params := authdb.CreateSessionParams{
		ID:         pgutil.UUID(s.ID()),
		UserID:     pgutil.UUID(s.SubjectID()),
		RefreshJti: pgutil.UUID(s.RefreshID()),
		ExpiresAt:  pgutil.Timestamptz(s.ExpiresAt()),
		CreatedAt:  pgutil.Timestamptz(s.CreatedAt()),
		RevokedAt:  pgutil.OptionalTimestamptz(s.RevokedAt()),
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
	session, err := r.q.GetSessionByID(ctx, pgutil.UUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Session{}, ports.ErrNotFound{Name: "session"}
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

	return domain.RehydrateSession(
		parsedID,
		parsedSubID,
		parsedRefreshID,
		session.CreatedAt.Time,
		session.ExpiresAt.Time,
		pgutil.TimestamptzPtr(session.RevokedAt),
	), nil
}

func (r *SessionRepo) ValidateAndRotateRefresh(ctx context.Context,
	sessionID domain.SessionID,
	expectedCurrent domain.RefreshTokenID,
	newRefreshID domain.RefreshTokenID,
	newExp time.Time,
) (bool, error) {
	rows, err := r.q.RotateSessionRefreshID(ctx, authdb.RotateSessionRefreshIDParams{
		ID:           pgutil.UUID(sessionID),
		RefreshJti:   pgutil.UUID(newRefreshID),
		ExpiresAt:    pgutil.Timestamptz(newExp),
		RefreshJti_2: pgutil.UUID(expectedCurrent),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "auth_sessions_refresh_jti_key" {
				return false, ports.ErrAlreadyExists{Name: "refresh jti"}
			}
		}
		return false, err
	}

	if rows == 1 {
		return true, nil
	}

	_, err = r.q.GetSessionByID(ctx, pgutil.UUID(sessionID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return false, ports.ErrSessionRefreshMismatch
}

func (r *SessionRepo) Revoke(ctx context.Context, sessionID domain.SessionID, now time.Time) error {
	err := r.q.RevokeSession(ctx, authdb.RevokeSessionParams{
		ID:        pgutil.UUID(sessionID),
		RevokedAt: pgutil.Timestamptz(now),
	})
	if err != nil {
		return err
	}

	return nil
}
