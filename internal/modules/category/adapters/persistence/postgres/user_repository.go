package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	categorydb "github.com/zchelalo/expense-control-back/internal/db/sqlc/category"
	"github.com/zchelalo/expense-control-back/internal/modules/category/domain"
	pgutil "github.com/zchelalo/expense-control-back/internal/shared/postgresutil"
)

type UserRepo struct {
	q *categorydb.Queries
}

func NewUserRepo(db categorydb.DBTX) *UserRepo {
	return &UserRepo{q: categorydb.New(db)}
}

func (r *UserRepo) Exists(ctx context.Context, userID domain.UserID) (bool, error) {
	exists, err := r.q.UserExists(ctx, pgutil.UUID(userID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return exists, nil
}
