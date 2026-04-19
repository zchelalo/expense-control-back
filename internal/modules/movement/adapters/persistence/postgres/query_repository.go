package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	movementdb "github.com/zchelalo/expense-control-back/internal/db/sqlc/movement"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/domain"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/ports"
)

type QueryRepo struct {
	q *movementdb.Queries
}

func NewQueryRepo(db movementdb.DBTX) *QueryRepo {
	return &QueryRepo{q: movementdb.New(db)}
}

func (r *QueryRepo) ByIDForUser(ctx context.Context, movementID domain.MovementID, userID domain.UserID) (domain.MovementDetails, error) {
	row, err := r.q.GetMovementDetailsByIDForUser(ctx, movementdb.GetMovementDetailsByIDForUserParams{
		ID:     toPgUUID(movementID),
		UserID: toPgUUID(userID),
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
	user := toPgUUID(userID)
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
			}
		})
	case 1:
		if filter.AccountID != nil {
			rows, err := r.q.ListMovementsByUserIDAndAccountIDAfter(ctx, movementdb.ListMovementsByUserIDAndAccountIDAfterParams{
				AccountID:        toPgUUID(*filter.AccountID),
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
				}
			})
		}

		if filter.CategoryID != nil {
			rows, err := r.q.ListMovementsByUserIDAndCategoryIDAfter(ctx, movementdb.ListMovementsByUserIDAndCategoryIDAfterParams{
				UserID:           user,
				CategoryID:       toPgUUID(*filter.CategoryID),
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
				}
			})
		}

		rows, err := r.q.ListMovementsByUserIDAndMovementTypeIDAfter(ctx, movementdb.ListMovementsByUserIDAndMovementTypeIDAfterParams{
			UserID:           user,
			MovementTypeID:   toPgUUID(*filter.MovementTypeID),
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
			}
		})
	default:
		rows, err := r.q.ListMovementsByUserIDFilteredAfter(ctx, movementdb.ListMovementsByUserIDFilteredAfterParams{
			UserID:           user,
			AccountID:        optionalAccountID(filter.AccountID),
			CategoryID:       optionalCategoryID(filter.CategoryID),
			MovementTypeID:   optionalMovementTypeID(filter.MovementTypeID),
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
			}
		})
	}
}

func (r *QueryRepo) listBefore(ctx context.Context, userID domain.UserID, filter ports.ListMovementsFilter) ([]domain.MovementDetails, error) {
	cursorCreatedAt, cursorMovementID := buildBeforeCursor(filter)
	user := toPgUUID(userID)
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
			}
		})
	case 1:
		switch {
		case filter.AccountID != nil:
			rows, queryErr := r.q.ListMovementsByUserIDAndAccountIDBefore(ctx, movementdb.ListMovementsByUserIDAndAccountIDBeforeParams{
				AccountID:        toPgUUID(*filter.AccountID),
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
				}
			})
		case filter.CategoryID != nil:
			rows, queryErr := r.q.ListMovementsByUserIDAndCategoryIDBefore(ctx, movementdb.ListMovementsByUserIDAndCategoryIDBeforeParams{
				UserID:           user,
				CategoryID:       toPgUUID(*filter.CategoryID),
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
				}
			})
		default:
			rows, queryErr := r.q.ListMovementsByUserIDAndMovementTypeIDBefore(ctx, movementdb.ListMovementsByUserIDAndMovementTypeIDBeforeParams{
				UserID:           user,
				MovementTypeID:   toPgUUID(*filter.MovementTypeID),
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
				}
			})
		}
	default:
		rows, queryErr := r.q.ListMovementsByUserIDFilteredBefore(ctx, movementdb.ListMovementsByUserIDFilteredBeforeParams{
			UserID:           user,
			AccountID:        optionalAccountID(filter.AccountID),
			CategoryID:       optionalCategoryID(filter.CategoryID),
			MovementTypeID:   optionalMovementTypeID(filter.MovementTypeID),
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
	var createdAt pgtype.Timestamptz
	if filter.CreatedAt != nil {
		createdAt = toPgTimestamptz(*filter.CreatedAt)
	}

	var movementID pgtype.UUID
	if filter.MovementID != nil {
		movementID = toPgUUID(*filter.MovementID)
	}

	return createdAt, movementID
}

func buildBeforeCursor(filter ports.ListMovementsFilter) (pgtype.Timestamptz, pgtype.UUID) {
	var createdAt pgtype.Timestamptz
	if filter.CreatedAt != nil {
		createdAt = toPgTimestamptz(*filter.CreatedAt)
	}

	var movementID pgtype.UUID
	if filter.MovementID != nil {
		movementID = toPgUUID(*filter.MovementID)
	}

	return createdAt, movementID
}

func optionalAccountID(accountID *domain.AccountID) pgtype.UUID {
	if accountID == nil {
		return pgtype.UUID{Valid: false}
	}

	return toPgUUID(*accountID)
}

func optionalCategoryID(categoryID *domain.CategoryID) pgtype.UUID {
	if categoryID == nil {
		return pgtype.UUID{Valid: false}
	}

	return toPgUUID(*categoryID)
}

func optionalMovementTypeID(movementTypeID *domain.MovementTypeID) pgtype.UUID {
	if movementTypeID == nil {
		return pgtype.UUID{Valid: false}
	}

	return toPgUUID(*movementTypeID)
}
