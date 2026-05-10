package list

import (
	"context"

	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/domain"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/ports"
	"go.uber.org/zap"
)

type UseCase struct {
	query           ports.QueryRepository
	users           ports.UserReferenceRepository
	pagLimitDefault int
}

func New(
	query ports.QueryRepository,
	users ports.UserReferenceRepository,
	pagLimitDefault int,
) *UseCase {
	return &UseCase{
		query:           query,
		users:           users,
		pagLimitDefault: pagLimitDefault,
	}
}

type Result struct {
	Movements []domain.MovementDetails
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (Result, error) {
	log := middleware.LoggerFrom(ctx)

	userID, err := domain.NewUserID(cmd.UserID)
	if err != nil {
		log.Warn("invalid user ID in list movements request",
			zap.String("stage", "validate_input"),
			zap.String("user_id", cmd.UserID.String()),
			zap.Error(err),
		)
		return Result{}, err
	}

	var accountID *domain.AccountID
	if cmd.AccountID != nil {
		id, err := domain.NewAccountID(*cmd.AccountID)
		if err != nil {
			log.Warn("invalid account ID in list movements request",
				zap.String("stage", "validate_input"),
				zap.String("account_id", cmd.AccountID.String()),
				zap.Error(err),
			)
			return Result{}, err
		}
		accountID = &id
	}

	var categoryID *domain.CategoryID
	if cmd.CategoryID != nil {
		id, err := domain.NewCategoryID(*cmd.CategoryID)
		if err != nil {
			log.Warn("invalid category ID in list movements request",
				zap.String("stage", "validate_input"),
				zap.String("category_id", cmd.CategoryID.String()),
				zap.Error(err),
			)
			return Result{}, err
		}
		categoryID = &id
	}

	var movementTypeID *domain.MovementTypeID
	if cmd.MovementTypeID != nil {
		id, err := domain.NewMovementTypeID(*cmd.MovementTypeID)
		if err != nil {
			log.Warn("invalid movement type ID in list movements request",
				zap.String("stage", "validate_input"),
				zap.String("movement_type_id", cmd.MovementTypeID.String()),
				zap.Error(err),
			)
			return Result{}, err
		}
		movementTypeID = &id
	}

	var movementID *domain.MovementID
	if cmd.MovementID != nil {
		id, err := domain.NewMovementID(*cmd.MovementID)
		if err != nil {
			log.Warn("invalid movement ID in list movements request",
				zap.String("stage", "validate_input"),
				zap.String("movement_id", cmd.MovementID.String()),
				zap.Error(err),
			)
			return Result{}, err
		}
		movementID = &id
	}

	if cmd.Limit <= 0 {
		log.Warn("invalid pagination limit in list movements request",
			zap.String("stage", "validate_input"),
			zap.Int("limit", cmd.Limit),
		)
		cmd.Limit = uc.pagLimitDefault
	}

	if cmd.CreatedAt == nil && movementID != nil {
		log.Warn("movement ID provided without created at in list movements request",
			zap.String("stage", "validate_input"),
			zap.String("movement_id", movementID.String()),
		)
		return Result{}, ErrMovementIDWithoutCreatedAt
	}

	if cmd.CreatedAt != nil && movementID == nil {
		log.Warn("created at provided without movement ID in list movements request",
			zap.String("stage", "validate_input"),
			zap.Time("created_at", *cmd.CreatedAt),
		)
		return Result{}, ErrCreatedAtWithoutMovementID
	}

	if cmd.DateFrom != nil && cmd.DateTo != nil && cmd.DateFrom.After(*cmd.DateTo) {
		log.Warn("invalid date range in list movements request",
			zap.String("stage", "validate_input"),
			zap.Time("date_from", *cmd.DateFrom),
			zap.Time("date_to", *cmd.DateTo),
		)
		return Result{}, ErrDateFromAfterDateTo
	}

	exists, err := uc.users.Exists(ctx, userID)
	if err != nil {
		log.Error("failed to check if user exists",
			zap.String("stage", "check_user_exists"),
			zap.Error(err),
		)
		return Result{}, err
	}
	if !exists {
		log.Warn("user not found for list movements request",
			zap.String("stage", "check_user_exists"),
			zap.String("user_id", userID.String()),
		)
		return Result{}, ports.ErrNotFound{Name: "user"}
	}

	movements, err := uc.query.ListByUserID(ctx, userID, ports.ListMovementsFilter{
		AccountID:      accountID,
		CategoryID:     categoryID,
		MovementTypeID: movementTypeID,
		DateFrom:       cmd.DateFrom,
		DateTo:         cmd.DateTo,
		CreatedAt:      cmd.CreatedAt,
		MovementID:     movementID,
		Limit:          cmd.Limit,
		IsBefore:       cmd.IsBefore,
	})
	if err != nil {
		log.Error("failed to list movements",
			zap.String("stage", "list_movements"),
			zap.Error(err),
		)
		return Result{}, err
	}

	return Result{Movements: movements}, nil
}
