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

type CredentialRepo struct {
  q *authdb.Queries
}

func NewCredentialRepo(db authdb.DBTX) *CredentialRepo {
  return &CredentialRepo{q: authdb.New(db)}
}

func (r *CredentialRepo) CreateAccount(ctx context.Context, acc domain.Account) (domain.SubjectID, error) {
	var deletedAt pgtype.Timestamptz
	if acc.DeletedAt() != nil {
		deletedAt = pgtype.Timestamptz{Time: *acc.DeletedAt(), Valid: true}
	} else {
		deletedAt = pgtype.Timestamptz{Valid: false}
	}
	params := authdb.CreateUserParams{
		ID: pgtype.UUID{
			Bytes: acc.ID().UUID(),
			Valid: true,
		},
		Email:     acc.Email().String(),
		Password:  acc.PasswordHash().String(),
		CreatedAt: pgtype.Timestamptz{
			Time: acc.CreatedAt(),
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamptz{
			Time: acc.UpdatedAt(),
			Valid: true,
		},
		DeletedAt: deletedAt,
	}

	subID, err := r.q.CreateUser(ctx, params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "users_email_key" {
				return domain.SubjectID{}, ports.ErrAlreadyExists{Name: "email"}
			}
		}
		return domain.SubjectID{}, err
	}

	parsedSubID, err := domain.NewSubjectID(subID.Bytes)
	if err != nil {
		return domain.SubjectID{}, err
	}

	return parsedSubID, nil
}

func (r *CredentialRepo) ByEmail(ctx context.Context, email domain.Email) (domain.Account, error) {
	user, err := r.q.GetUserByEmail(ctx, email.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Account{}, ports.ErrNotFound{Name:"user"}
		}
		return domain.Account{}, err
	}

	parsedSubID, err := domain.NewSubjectID(user.ID.Bytes)
	if err != nil {
		return domain.Account{}, err
	}
	parsedEmail, err := domain.NewEmail(user.Email)
	if err != nil {
		return domain.Account{}, err
	}
	parsedPasswordHash, err := domain.NewPasswordHashFromHash(user.Password)
	if err != nil {
		return domain.Account{}, err
	}
	var parsedDeletedAt *time.Time
	if user.DeletedAt.Valid {
		t := user.DeletedAt.Time
		parsedDeletedAt = &t
	}
	account := domain.RehydrateAccount(
		parsedSubID,
    parsedEmail,
		parsedPasswordHash,
		user.CreatedAt.Time,
		user.UpdatedAt.Time,
		parsedDeletedAt,
	)

	return account, nil
}