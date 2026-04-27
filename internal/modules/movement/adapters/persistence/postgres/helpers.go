package postgres

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/domain"
	"github.com/zchelalo/expense-control-back/internal/shared/postgresutil"
)

type movementDetailsFields struct {
	ID               pgtype.UUID
	Amount           pgtype.Numeric
	Description      string
	MovementTypeID   pgtype.UUID
	CategoryID       pgtype.UUID
	AccountID        pgtype.UUID
	UserID           pgtype.UUID
	CreatedAt        pgtype.Timestamptz
	UpdatedAt        pgtype.Timestamptz
	DeletedAt        pgtype.Timestamptz
	MovementTypeKey  string
	MovementTypeName string
	CategoryName     string
	AccountName      string
}

func hydrateMovement(
	id pgtype.UUID,
	amount pgtype.Numeric,
	description string,
	movementTypeID pgtype.UUID,
	categoryID pgtype.UUID,
	accountID pgtype.UUID,
	userID pgtype.UUID,
	createdAt pgtype.Timestamptz,
	updatedAt pgtype.Timestamptz,
	deletedAt pgtype.Timestamptz,
) (domain.Movement, error) {
	parsedID, err := domain.NewMovementID(id.Bytes)
	if err != nil {
		return domain.Movement{}, err
	}

	parsedAmountValue, err := postgresutil.NumericToFloat64(amount)
	if err != nil {
		return domain.Movement{}, err
	}

	parsedAmount, err := domain.NewAmount(parsedAmountValue)
	if err != nil {
		return domain.Movement{}, err
	}

	parsedDescription, err := domain.NewDescription(description)
	if err != nil {
		return domain.Movement{}, err
	}

	parsedMovementTypeID, err := domain.NewMovementTypeID(movementTypeID.Bytes)
	if err != nil {
		return domain.Movement{}, err
	}

	parsedCategoryID, err := domain.NewCategoryID(categoryID.Bytes)
	if err != nil {
		return domain.Movement{}, err
	}

	parsedAccountID, err := domain.NewAccountID(accountID.Bytes)
	if err != nil {
		return domain.Movement{}, err
	}

	parsedUserID, err := domain.NewUserID(userID.Bytes)
	if err != nil {
		return domain.Movement{}, err
	}

	return domain.RehydrateMovement(
		parsedID,
		parsedAmount,
		parsedDescription,
		parsedMovementTypeID,
		parsedCategoryID,
		parsedAccountID,
		parsedUserID,
		createdAt.Time,
		updatedAt.Time,
		postgresutil.TimestamptzPtr(deletedAt),
	), nil
}

func hydrateMovementDetails(fields movementDetailsFields) (domain.MovementDetails, error) {
	movement, err := hydrateMovement(
		fields.ID,
		fields.Amount,
		fields.Description,
		fields.MovementTypeID,
		fields.CategoryID,
		fields.AccountID,
		fields.UserID,
		fields.CreatedAt,
		fields.UpdatedAt,
		fields.DeletedAt,
	)
	if err != nil {
		return domain.MovementDetails{}, err
	}

	movementType, err := domain.RehydrateMovementType(
		movement.MovementTypeID(),
		fields.MovementTypeKey,
		fields.MovementTypeName,
	)
	if err != nil {
		return domain.MovementDetails{}, err
	}

	category, err := domain.RehydrateCategory(
		movement.CategoryID(),
		fields.CategoryName,
	)
	if err != nil {
		return domain.MovementDetails{}, err
	}

	account, err := domain.RehydrateAccount(
		movement.AccountID(),
		fields.AccountName,
	)
	if err != nil {
		return domain.MovementDetails{}, err
	}

	return domain.NewMovementDetails(movement, movementType, category, account), nil
}

func mapMovementDetailsRows[T any](rows []T, extract func(T) movementDetailsFields) ([]domain.MovementDetails, error) {
	result := make([]domain.MovementDetails, len(rows))

	for i, row := range rows {
		details, err := hydrateMovementDetails(extract(row))
		if err != nil {
			return nil, err
		}

		result[i] = details
	}

	return result, nil
}

func reverseMovementDetails(items []domain.MovementDetails) {
	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}
}
