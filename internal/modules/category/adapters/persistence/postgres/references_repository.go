package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	categorydb "github.com/zchelalo/expense-control-back/internal/db/sqlc/category"
	"github.com/zchelalo/expense-control-back/internal/modules/category/domain"
	pgutil "github.com/zchelalo/expense-control-back/internal/shared/postgresutil"
)

type UserReferenceRepository struct {
	q *categorydb.Queries
}

func NewUserReferenceRepository(db categorydb.DBTX) *UserReferenceRepository {
	return &UserReferenceRepository{q: categorydb.New(db)}
}

func (r *UserReferenceRepository) Exists(ctx context.Context, userID domain.UserID) (bool, error) {
	exists, err := r.q.UserExists(ctx, pgutil.UUID(userID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return exists, nil
}
