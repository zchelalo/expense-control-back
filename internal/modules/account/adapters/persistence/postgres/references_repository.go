package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	accountdb "github.com/zchelalo/expense-control-back/internal/db/sqlc/account"
	"github.com/zchelalo/expense-control-back/internal/modules/account/domain"
	pgutil "github.com/zchelalo/expense-control-back/internal/shared/postgresutil"
)

type UserReferenceRepository struct {
	q *accountdb.Queries
}

func NewUserReferenceRepository(db accountdb.DBTX) *UserReferenceRepository {
	return &UserReferenceRepository{q: accountdb.New(db)}
}

func (r *UserReferenceRepository) Exists(ctx context.Context, id domain.UserID) (bool, error) {
	exists, err := r.q.UserExists(ctx, pgutil.UUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return exists, nil
}
