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

type UserReferenceRepository struct {
	q *movementdb.Queries
}

func NewUserReferenceRepository(db movementdb.DBTX) *UserReferenceRepository {
	return &UserReferenceRepository{q: movementdb.New(db)}
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

type AccountReferenceRepository struct {
	q *movementdb.Queries
}

func NewAccountReferenceRepository(db movementdb.DBTX) *AccountReferenceRepository {
	return &AccountReferenceRepository{q: movementdb.New(db)}
}

func (r *AccountReferenceRepository) ExistsByUserID(ctx context.Context, accountID domain.AccountID, userID domain.UserID) (bool, error) {
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

type MovementTypeReferenceRepository struct {
	q *movementdb.Queries
}

func NewMovementTypeReferenceRepository(db movementdb.DBTX) *MovementTypeReferenceRepository {
	return &MovementTypeReferenceRepository{q: movementdb.New(db)}
}

func (r *MovementTypeReferenceRepository) ByID(ctx context.Context, movementTypeID domain.MovementTypeID) (domain.MovementType, error) {
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

type CategoryReferenceRepository struct {
	q *movementdb.Queries
}

func NewCategoryReferenceRepository(db movementdb.DBTX) *CategoryReferenceRepository {
	return &CategoryReferenceRepository{q: movementdb.New(db)}
}

func (r *CategoryReferenceRepository) ByIDForUser(ctx context.Context, categoryID domain.CategoryID, userID domain.UserID) (domain.Category, error) {
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
