package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	movementdb "github.com/zchelalo/expense-control-back/internal/db/sqlc/movement"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/domain"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/ports"
	pgutil "github.com/zchelalo/expense-control-back/internal/shared/postgresutil"
)

type QueryRepo struct {
	q *movementdb.Queries
}

func NewQueryRepo(db movementdb.DBTX) *QueryRepo {
	return &QueryRepo{q: movementdb.New(db)}
}

func (r *QueryRepo) ByIDForUser(ctx context.Context, movementID domain.MovementID, userID domain.UserID) (domain.MovementDetails, error) {
	row, err := r.q.GetMovementDetailsByIDForUser(ctx, movementdb.GetMovementDetailsByIDForUserParams{
		ID:     pgutil.UUID(movementID),
		UserID: pgutil.UUID(userID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.MovementDetails{}, ports.ErrNotFound{Name: "movement"}
		}
		return domain.MovementDetails{}, err
	}

	return hydrateMovementDetails(movementDetailsFields{
		ID:               row.ID,
		Amount:           row.Amount,
		Description:      row.Description,
		MovementTypeID:   row.MovementTypeID,
		CategoryID:       row.CategoryID,
		AccountID:        row.AccountID,
		UserID:           row.UserID,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
		DeletedAt:        row.DeletedAt,
		MovementTypeKey:  row.MovementTypeKey,
		MovementTypeName: row.MovementTypeName,
		CategoryName:     row.CategoryName,
		AccountName:      row.AccountName,
	})
}

func (r *QueryRepo) ListByUserID(ctx context.Context, userID domain.UserID, filter ports.ListMovementsFilter) ([]domain.MovementDetails, error) {
	if filter.IsBefore {
		return r.listBefore(ctx, userID, filter)
	}

	return r.listAfter(ctx, userID, filter)
}

func (r *QueryRepo) listAfter(ctx context.Context, userID domain.UserID, filter ports.ListMovementsFilter) ([]domain.MovementDetails, error) {
	cursorCreatedAt, cursorMovementID := buildAfterCursor(filter)
	user := pgutil.UUID(userID)
	limit := int32(filter.Limit)

	switch countActiveFilters(filter) {
	case 0:
		rows, err := r.q.ListMovementsByUserIDAfter(ctx, movementdb.ListMovementsByUserIDAfterParams{
			UserID:           user,
			CursorCreatedAt:  cursorCreatedAt,
			CursorMovementID: cursorMovementID,
			LimitCount:       limit,
		})
		if err != nil {
			return nil, err
		}

		return mapMovementDetailsRows(rows, func(row movementdb.ListMovementsByUserIDAfterRow) movementDetailsFields {
			return movementDetailsFields{
				ID:               row.ID,
				Amount:           row.Amount,
				Description:      row.Description,
				MovementTypeID:   row.MovementTypeID,
				CategoryID:       row.CategoryID,
				AccountID:        row.AccountID,
				UserID:           row.UserID,
				CreatedAt:        row.CreatedAt,
				UpdatedAt:        row.UpdatedAt,
				DeletedAt:        row.DeletedAt,
				MovementTypeKey:  row.MovementTypeKey,
				MovementTypeName: row.MovementTypeName,
				CategoryName:     row.CategoryName,
				AccountName:      row.AccountName,
			}
		})
	case 1:
		if filter.AccountID != nil {
			rows, err := r.q.ListMovementsByUserIDAndAccountIDAfter(ctx, movementdb.ListMovementsByUserIDAndAccountIDAfterParams{
				AccountID:        pgutil.UUID(*filter.AccountID),
				UserID:           user,
				CursorCreatedAt:  cursorCreatedAt,
				CursorMovementID: cursorMovementID,
				LimitCount:       limit,
			})
			if err != nil {
				return nil, err
			}

			return mapMovementDetailsRows(rows, func(row movementdb.ListMovementsByUserIDAndAccountIDAfterRow) movementDetailsFields {
				return movementDetailsFields{
					ID:               row.ID,
					Amount:           row.Amount,
					Description:      row.Description,
					MovementTypeID:   row.MovementTypeID,
					CategoryID:       row.CategoryID,
					AccountID:        row.AccountID,
					UserID:           row.UserID,
					CreatedAt:        row.CreatedAt,
					UpdatedAt:        row.UpdatedAt,
					DeletedAt:        row.DeletedAt,
					MovementTypeKey:  row.MovementTypeKey,
					MovementTypeName: row.MovementTypeName,
					CategoryName:     row.CategoryName,
					AccountName:      row.AccountName,
				}
			})
		}

		if filter.CategoryID != nil {
			rows, err := r.q.ListMovementsByUserIDAndCategoryIDAfter(ctx, movementdb.ListMovementsByUserIDAndCategoryIDAfterParams{
				UserID:           user,
				CategoryID:       pgutil.UUID(*filter.CategoryID),
				CursorCreatedAt:  cursorCreatedAt,
				CursorMovementID: cursorMovementID,
				LimitCount:       limit,
			})
			if err != nil {
				return nil, err
			}

			return mapMovementDetailsRows(rows, func(row movementdb.ListMovementsByUserIDAndCategoryIDAfterRow) movementDetailsFields {
				return movementDetailsFields{
					ID:               row.ID,
					Amount:           row.Amount,
					Description:      row.Description,
					MovementTypeID:   row.MovementTypeID,
					CategoryID:       row.CategoryID,
					AccountID:        row.AccountID,
					UserID:           row.UserID,
					CreatedAt:        row.CreatedAt,
					UpdatedAt:        row.UpdatedAt,
					DeletedAt:        row.DeletedAt,
					MovementTypeKey:  row.MovementTypeKey,
					MovementTypeName: row.MovementTypeName,
					CategoryName:     row.CategoryName,
					AccountName:      row.AccountName,
				}
			})
		}

		rows, err := r.q.ListMovementsByUserIDAndMovementTypeIDAfter(ctx, movementdb.ListMovementsByUserIDAndMovementTypeIDAfterParams{
			UserID:           user,
			MovementTypeID:   pgutil.UUID(*filter.MovementTypeID),
			CursorCreatedAt:  cursorCreatedAt,
			CursorMovementID: cursorMovementID,
			LimitCount:       limit,
		})
		if err != nil {
			return nil, err
		}

		return mapMovementDetailsRows(rows, func(row movementdb.ListMovementsByUserIDAndMovementTypeIDAfterRow) movementDetailsFields {
			return movementDetailsFields{
				ID:               row.ID,
				Amount:           row.Amount,
				Description:      row.Description,
				MovementTypeID:   row.MovementTypeID,
				CategoryID:       row.CategoryID,
				AccountID:        row.AccountID,
				UserID:           row.UserID,
				CreatedAt:        row.CreatedAt,
				UpdatedAt:        row.UpdatedAt,
				DeletedAt:        row.DeletedAt,
				MovementTypeKey:  row.MovementTypeKey,
				MovementTypeName: row.MovementTypeName,
				CategoryName:     row.CategoryName,
				AccountName:      row.AccountName,
			}
		})
	default:
		rows, err := r.q.ListMovementsByUserIDFilteredAfter(ctx, movementdb.ListMovementsByUserIDFilteredAfterParams{
			UserID:           user,
			AccountID:        pgutil.OptionalUUID(filter.AccountID),
			CategoryID:       pgutil.OptionalUUID(filter.CategoryID),
			MovementTypeID:   pgutil.OptionalUUID(filter.MovementTypeID),
			CursorCreatedAt:  cursorCreatedAt,
			CursorMovementID: cursorMovementID,
			LimitCount:       limit,
		})
		if err != nil {
			return nil, err
		}

		return mapMovementDetailsRows(rows, func(row movementdb.ListMovementsByUserIDFilteredAfterRow) movementDetailsFields {
			return movementDetailsFields{
				ID:               row.ID,
				Amount:           row.Amount,
				Description:      row.Description,
				MovementTypeID:   row.MovementTypeID,
				CategoryID:       row.CategoryID,
				AccountID:        row.AccountID,
				UserID:           row.UserID,
				CreatedAt:        row.CreatedAt,
				UpdatedAt:        row.UpdatedAt,
				DeletedAt:        row.DeletedAt,
				MovementTypeKey:  row.MovementTypeKey,
				MovementTypeName: row.MovementTypeName,
				CategoryName:     row.CategoryName,
				AccountName:      row.AccountName,
			}
		})
	}
}

func (r *QueryRepo) listBefore(ctx context.Context, userID domain.UserID, filter ports.ListMovementsFilter) ([]domain.MovementDetails, error) {
	cursorCreatedAt, cursorMovementID := buildBeforeCursor(filter)
	user := pgutil.UUID(userID)
	limit := int32(filter.Limit)

	var (
		result []domain.MovementDetails
		err    error
	)

	switch countActiveFilters(filter) {
	case 0:
		rows, queryErr := r.q.ListMovementsByUserIDBefore(ctx, movementdb.ListMovementsByUserIDBeforeParams{
			UserID:           user,
			CursorCreatedAt:  cursorCreatedAt,
			CursorMovementID: cursorMovementID,
			LimitCount:       limit,
		})
		if queryErr != nil {
			return nil, queryErr
		}

		result, err = mapMovementDetailsRows(rows, func(row movementdb.ListMovementsByUserIDBeforeRow) movementDetailsFields {
			return movementDetailsFields{
				ID:               row.ID,
				Amount:           row.Amount,
				Description:      row.Description,
				MovementTypeID:   row.MovementTypeID,
				CategoryID:       row.CategoryID,
				AccountID:        row.AccountID,
				UserID:           row.UserID,
				CreatedAt:        row.CreatedAt,
				UpdatedAt:        row.UpdatedAt,
				DeletedAt:        row.DeletedAt,
				MovementTypeKey:  row.MovementTypeKey,
				MovementTypeName: row.MovementTypeName,
				CategoryName:     row.CategoryName,
				AccountName:      row.AccountName,
			}
		})
	case 1:
		switch {
		case filter.AccountID != nil:
			rows, queryErr := r.q.ListMovementsByUserIDAndAccountIDBefore(ctx, movementdb.ListMovementsByUserIDAndAccountIDBeforeParams{
				AccountID:        pgutil.UUID(*filter.AccountID),
				UserID:           user,
				CursorCreatedAt:  cursorCreatedAt,
				CursorMovementID: cursorMovementID,
				LimitCount:       limit,
			})
			if queryErr != nil {
				return nil, queryErr
			}

			result, err = mapMovementDetailsRows(rows, func(row movementdb.ListMovementsByUserIDAndAccountIDBeforeRow) movementDetailsFields {
				return movementDetailsFields{
					ID:               row.ID,
					Amount:           row.Amount,
					Description:      row.Description,
					MovementTypeID:   row.MovementTypeID,
					CategoryID:       row.CategoryID,
					AccountID:        row.AccountID,
					UserID:           row.UserID,
					CreatedAt:        row.CreatedAt,
					UpdatedAt:        row.UpdatedAt,
					DeletedAt:        row.DeletedAt,
					MovementTypeKey:  row.MovementTypeKey,
					MovementTypeName: row.MovementTypeName,
					CategoryName:     row.CategoryName,
					AccountName:      row.AccountName,
				}
			})
		case filter.CategoryID != nil:
			rows, queryErr := r.q.ListMovementsByUserIDAndCategoryIDBefore(ctx, movementdb.ListMovementsByUserIDAndCategoryIDBeforeParams{
				UserID:           user,
				CategoryID:       pgutil.UUID(*filter.CategoryID),
				CursorCreatedAt:  cursorCreatedAt,
				CursorMovementID: cursorMovementID,
				LimitCount:       limit,
			})
			if queryErr != nil {
				return nil, queryErr
			}

			result, err = mapMovementDetailsRows(rows, func(row movementdb.ListMovementsByUserIDAndCategoryIDBeforeRow) movementDetailsFields {
				return movementDetailsFields{
					ID:               row.ID,
					Amount:           row.Amount,
					Description:      row.Description,
					MovementTypeID:   row.MovementTypeID,
					CategoryID:       row.CategoryID,
					AccountID:        row.AccountID,
					UserID:           row.UserID,
					CreatedAt:        row.CreatedAt,
					UpdatedAt:        row.UpdatedAt,
					DeletedAt:        row.DeletedAt,
					MovementTypeKey:  row.MovementTypeKey,
					MovementTypeName: row.MovementTypeName,
					CategoryName:     row.CategoryName,
					AccountName:      row.AccountName,
				}
			})
		default:
			rows, queryErr := r.q.ListMovementsByUserIDAndMovementTypeIDBefore(ctx, movementdb.ListMovementsByUserIDAndMovementTypeIDBeforeParams{
				UserID:           user,
				MovementTypeID:   pgutil.UUID(*filter.MovementTypeID),
				CursorCreatedAt:  cursorCreatedAt,
				CursorMovementID: cursorMovementID,
				LimitCount:       limit,
			})
			if queryErr != nil {
				return nil, queryErr
			}

			result, err = mapMovementDetailsRows(rows, func(row movementdb.ListMovementsByUserIDAndMovementTypeIDBeforeRow) movementDetailsFields {
				return movementDetailsFields{
					ID:               row.ID,
					Amount:           row.Amount,
					Description:      row.Description,
					MovementTypeID:   row.MovementTypeID,
					CategoryID:       row.CategoryID,
					AccountID:        row.AccountID,
					UserID:           row.UserID,
					CreatedAt:        row.CreatedAt,
					UpdatedAt:        row.UpdatedAt,
					DeletedAt:        row.DeletedAt,
					MovementTypeKey:  row.MovementTypeKey,
					MovementTypeName: row.MovementTypeName,
					CategoryName:     row.CategoryName,
					AccountName:      row.AccountName,
				}
			})
		}
	default:
		rows, queryErr := r.q.ListMovementsByUserIDFilteredBefore(ctx, movementdb.ListMovementsByUserIDFilteredBeforeParams{
			UserID:           user,
			AccountID:        pgutil.OptionalUUID(filter.AccountID),
			CategoryID:       pgutil.OptionalUUID(filter.CategoryID),
			MovementTypeID:   pgutil.OptionalUUID(filter.MovementTypeID),
			CursorCreatedAt:  cursorCreatedAt,
			CursorMovementID: cursorMovementID,
			LimitCount:       limit,
		})
		if queryErr != nil {
			return nil, queryErr
		}

		result, err = mapMovementDetailsRows(rows, func(row movementdb.ListMovementsByUserIDFilteredBeforeRow) movementDetailsFields {
			return movementDetailsFields{
				ID:               row.ID,
				Amount:           row.Amount,
				Description:      row.Description,
				MovementTypeID:   row.MovementTypeID,
				CategoryID:       row.CategoryID,
				AccountID:        row.AccountID,
				UserID:           row.UserID,
				CreatedAt:        row.CreatedAt,
				UpdatedAt:        row.UpdatedAt,
				DeletedAt:        row.DeletedAt,
				MovementTypeKey:  row.MovementTypeKey,
				MovementTypeName: row.MovementTypeName,
				CategoryName:     row.CategoryName,
				AccountName:      row.AccountName,
			}
		})
	}

	if err != nil {
		return nil, err
	}

	reverseMovementDetails(result)
	return result, nil
}

func countActiveFilters(filter ports.ListMovementsFilter) int {
	count := 0
	if filter.AccountID != nil {
		count++
	}
	if filter.CategoryID != nil {
		count++
	}
	if filter.MovementTypeID != nil {
		count++
	}

	return count
}

func buildAfterCursor(filter ports.ListMovementsFilter) (pgtype.Timestamptz, pgtype.UUID) {
	return pgutil.OptionalTimestamptz(filter.CreatedAt), pgutil.OptionalUUID(filter.MovementID)
}

func buildBeforeCursor(filter ports.ListMovementsFilter) (pgtype.Timestamptz, pgtype.UUID) {
	return pgutil.OptionalTimestamptz(filter.CreatedAt), pgutil.OptionalUUID(filter.MovementID)
}
