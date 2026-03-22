package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	accountdb "github.com/zchelalo/expense-control-back/internal/db/sqlc/account"
	"github.com/zchelalo/expense-control-back/internal/modules/account/domain"
)

type UserRepo struct {
  q *accountdb.Queries
}

func NewUserRepo(db accountdb.DBTX) *UserRepo {
  return &UserRepo{q: accountdb.New(db)}
}

func (r *UserRepo) Exists(ctx context.Context, id domain.UserID) (bool, error) {
	exists, err := r.q.UserExists(ctx, pgtype.UUID{
		Bytes: id.UUID(),
		Valid: true,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return exists, nil
}