package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	movementdb "github.com/zchelalo/expense-control-back/internal/db/sqlc/movement"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/domain"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/ports"
	pgutil "github.com/zchelalo/expense-control-back/internal/shared/postgresutil"
)

type UserRepo struct {
	q *movementdb.Queries
}

func NewUserRepo(db movementdb.DBTX) *UserRepo {
	return &UserRepo{q: movementdb.New(db)}
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

type AccountRepo struct {
	q *movementdb.Queries
}

func NewAccountRepo(db movementdb.DBTX) *AccountRepo {
	return &AccountRepo{q: movementdb.New(db)}
}

func (r *AccountRepo) ExistsByUserID(ctx context.Context, accountID domain.AccountID, userID domain.UserID) (bool, error) {
	exists, err := r.q.AccountExistsByUserID(ctx, movementdb.AccountExistsByUserIDParams{
		ID:     pgutil.UUID(accountID),
		UserID: pgutil.UUID(userID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return exists, nil
}

type MovementTypeRepo struct {
	q *movementdb.Queries
}

func NewMovementTypeRepo(db movementdb.DBTX) *MovementTypeRepo {
	return &MovementTypeRepo{q: movementdb.New(db)}
}

func (r *MovementTypeRepo) ByID(ctx context.Context, movementTypeID domain.MovementTypeID) (domain.MovementType, error) {
	movementType, err := r.q.GetMovementTypeByID(ctx, pgutil.UUID(movementTypeID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.MovementType{}, ports.ErrNotFound{Name: "movement type"}
		}
		return domain.MovementType{}, err
	}

	return domain.RehydrateMovementType(
		movementTypeID,
		movementType.Key,
		movementType.Name,
	)
}

type CategoryRepo struct {
	q *movementdb.Queries
}

func NewCategoryRepo(db movementdb.DBTX) *CategoryRepo {
	return &CategoryRepo{q: movementdb.New(db)}
}

func (r *CategoryRepo) ByIDForUser(ctx context.Context, categoryID domain.CategoryID, userID domain.UserID) (domain.Category, error) {
	category, err := r.q.GetCategoryByIDForUser(ctx, movementdb.GetCategoryByIDForUserParams{
		CategoryID: pgutil.UUID(categoryID),
		UserID:     pgutil.UUID(userID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Category{}, ports.ErrNotFound{Name: "category"}
		}
		return domain.Category{}, err
	}

	return domain.RehydrateCategory(categoryID, category.Name)
}
