package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	categorydb "github.com/zchelalo/expense-control-back/internal/db/sqlc/category"
	"github.com/zchelalo/expense-control-back/internal/modules/category/domain"
	"github.com/zchelalo/expense-control-back/internal/modules/category/ports"
	pgutil "github.com/zchelalo/expense-control-back/internal/shared/postgresutil"
)

type categoryTxDB interface {
	categorydb.DBTX
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type CategoryRepo struct {
	db categoryTxDB
	q  *categorydb.Queries
}

func NewCategoryRepo(db categoryTxDB) *CategoryRepo {
	return &CategoryRepo{
		db: db,
		q:  categorydb.New(db),
	}
}

func (r *CategoryRepo) Create(ctx context.Context, category domain.Category) (domain.Category, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.Category{}, err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)

	globalCategory, err := q.UpsertCategoryByName(ctx, categorydb.UpsertCategoryByNameParams{
		ID:        pgutil.UUID(category.ID()),
		Name:      category.Name().String(),
		CreatedAt: pgutil.Timestamptz(category.CreatedAt()),
		UpdatedAt: pgutil.Timestamptz(category.UpdatedAt()),
		DeletedAt: pgutil.OptionalTimestamptz(nil),
	})
	if err != nil {
		return domain.Category{}, err
	}

	actualCategoryID, err := domain.NewCategoryID(globalCategory.ID.Bytes)
	if err != nil {
		return domain.Category{}, err
	}

	userCategory, err := q.UpsertUserCategory(ctx, categorydb.UpsertUserCategoryParams{
		UserID:     pgutil.UUID(category.UserID()),
		CategoryID: pgutil.UUID(actualCategoryID),
		CreatedAt:  pgutil.Timestamptz(category.CreatedAt()),
		UpdatedAt:  pgutil.Timestamptz(category.UpdatedAt()),
		DeletedAt:  pgutil.OptionalTimestamptz(nil),
	})
	if err != nil {
		return domain.Category{}, err
	}

	createdCategory, err := domain.RehydrateCategory(
		actualCategoryID,
		globalCategory.Name,
		category.UserID(),
		userCategory.CreatedAt.Time,
		userCategory.UpdatedAt.Time,
	)
	if err != nil {
		return domain.Category{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Category{}, err
	}

	return createdCategory, nil
}

func (r *CategoryRepo) ListByUserID(ctx context.Context, userID domain.UserID, createdAt *time.Time, categoryID *domain.CategoryID, limit int, isBefore bool) ([]domain.Category, error) {
	if isBefore {
		rows, err := r.q.ListCategoriesByUserIDBefore(ctx, categorydb.ListCategoriesByUserIDBeforeParams{
			UserID:           pgutil.UUID(userID),
			CursorCreatedAt:  pgutil.OptionalTimestamptz(createdAt),
			CursorCategoryID: pgutil.OptionalUUID(categoryID),
			LimitCount:       int32(limit),
		})
		if err != nil {
			return nil, err
		}

		result := make([]domain.Category, len(rows))
		for i, row := range rows {
			item, err := hydrateCategory(row.ID, row.Name, row.UserID, row.CreatedAt, row.UpdatedAt)
			if err != nil {
				return nil, err
			}
			result[i] = item
		}

		reverseCategories(result)
		return result, nil
	}

	rows, err := r.q.ListCategoriesByUserIDAfter(ctx, categorydb.ListCategoriesByUserIDAfterParams{
		UserID:           pgutil.UUID(userID),
		CursorCreatedAt:  pgutil.OptionalTimestamptz(createdAt),
		CursorCategoryID: pgutil.OptionalUUID(categoryID),
		LimitCount:       int32(limit),
	})
	if err != nil {
		return nil, err
	}

	result := make([]domain.Category, len(rows))
	for i, row := range rows {
		item, err := hydrateCategory(row.ID, row.Name, row.UserID, row.CreatedAt, row.UpdatedAt)
		if err != nil {
			return nil, err
		}
		result[i] = item
	}

	return result, nil
}

func (r *CategoryRepo) Delete(ctx context.Context, userID domain.UserID, categoryID domain.CategoryID, now time.Time) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)

	affectedRows, err := q.DeleteUserCategory(ctx, categorydb.DeleteUserCategoryParams{
		UserID:     pgutil.UUID(userID),
		CategoryID: pgutil.UUID(categoryID),
		DeletedAt:  pgutil.Timestamptz(now),
		UpdatedAt:  pgutil.Timestamptz(now),
	})
	if err != nil {
		return err
	}
	if affectedRows > 0 {
		return tx.Commit(ctx)
	}

	userCategory, err := q.GetUserCategoryByUserIDAndCategoryID(ctx, categorydb.GetUserCategoryByUserIDAndCategoryIDParams{
		UserID:     pgutil.UUID(userID),
		CategoryID: pgutil.UUID(categoryID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ports.ErrNotFound{Name: "category"}
		}

		return err
	}
	if userCategory.DeletedAt.Valid {
		return ports.ErrNotFound{Name: "category"}
	}

	inUse, err := q.CategoryHasActiveMovementsForUser(ctx, categorydb.CategoryHasActiveMovementsForUserParams{
		UserID:     pgutil.UUID(userID),
		CategoryID: pgutil.UUID(categoryID),
	})
	if err != nil {
		return err
	}
	if inUse {
		return ports.ErrInUse{Name: "category"}
	}

	return ports.ErrNotFound{Name: "category"}
}

func hydrateCategory(
	categoryID pgtype.UUID,
	name string,
	userID pgtype.UUID,
	createdAt pgtype.Timestamptz,
	updatedAt pgtype.Timestamptz,
) (domain.Category, error) {
	parsedCategoryID, err := domain.NewCategoryID(categoryID.Bytes)
	if err != nil {
		return domain.Category{}, err
	}

	parsedUserID, err := domain.NewUserID(userID.Bytes)
	if err != nil {
		return domain.Category{}, err
	}

	return domain.RehydrateCategory(
		parsedCategoryID,
		name,
		parsedUserID,
		createdAt.Time,
		updatedAt.Time,
	)
}

func reverseCategories(items []domain.Category) {
	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}
}
