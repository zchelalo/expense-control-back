package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	authdb "github.com/zchelalo/expense-control-back/internal/db/sqlc/auth"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/domain"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/ports"
	pgutil "github.com/zchelalo/expense-control-back/internal/shared/postgresutil"
)

type CredentialRepo struct {
	q *authdb.Queries
}

func NewCredentialRepo(db authdb.DBTX) *CredentialRepo {
	return &CredentialRepo{q: authdb.New(db)}
}

func (r *CredentialRepo) CreateAccount(ctx context.Context, acc domain.Account) (domain.SubjectID, error) {
	params := authdb.CreateUserParams{
		ID:        pgutil.UUID(acc.ID()),
		Email:     acc.Email().String(),
		Password:  acc.PasswordHash().String(),
		CreatedAt: pgutil.Timestamptz(acc.CreatedAt()),
		UpdatedAt: pgutil.Timestamptz(acc.UpdatedAt()),
		DeletedAt: pgutil.OptionalTimestamptz(acc.DeletedAt()),
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
			return domain.Account{}, ports.ErrNotFound{Name: "user"}
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
	account := domain.RehydrateAccount(
		parsedSubID,
		parsedEmail,
		parsedPasswordHash,
		user.CreatedAt.Time,
		user.UpdatedAt.Time,
		pgutil.TimestamptzPtr(user.DeletedAt),
	)

	return account, nil
}
